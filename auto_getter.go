package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
)

type AutoProxyGetter interface {
	ProxyGetter
	CrawlProxy(url string) // Crawl available proxies from the Internet
}

//
type TimerAutoProxyGetter struct {
	DefaultProxyGetter
	interval int // How many seconds between crawling proxies
}

type MiniThresholdAutoProxyGetter struct {
	DefaultProxyGetter
	min int // Actively crawl when the proxy is lower than this value
}

func (p *TimerAutoProxyGetter) CrawlProxy(url string) {
	timeTickerChan := time.Tick(time.Second * time.Duration(p.interval))
	// init collyCollector
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36"),
		colly.AllowURLRevisit(),
	)
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	})
	for {
		if p.proxies != nil && len(p.proxies) != 0 {
			ipPort := []string{}
			for _, v := range p.proxies {
				ipPort = append(ipPort, "socks5://"+v)
			}
			rp, err := proxy.RoundRobinProxySwitcher(ipPort...)
			if err != nil {
				log.Fatal(err)
			}
			// Set proxy ip, here is polling
			c.SetProxyFunc(rp)
		}

		proxyList := []string{}
		c.OnRequest(func(r *colly.Request) {
			// fmt.Println("Visiting", r.URL)
		})
		c.OnError(func(_ *colly.Response, err error) {
			// fmt.Println("Something went wrong:", err)
		})
		c.OnResponse(func(r *colly.Response) {
			doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
			if err != nil {
				panic(err)
			}
			nodes := htmlquery.Find(doc, `//tbody/tr`)
			for _, node := range nodes {
				addr := htmlquery.FindOne(node, `./td[1]/text()`)
				port := htmlquery.FindOne(node, `./td[2]/text()`)
				proxy := htmlquery.InnerText(addr) + ":" + htmlquery.InnerText(port)
				proxyList = append(proxyList, proxy)
			}
			pageList := make(map[string]int)
			// nextPage
			pages := htmlquery.Find(doc, `//*[@id="listnav"]//a/@href`)
			for _, page := range pages {
				_, ok := pageList[htmlquery.InnerText(page)]
				if !ok {
					pageList[htmlquery.InnerText(page)] = 1
					go p.CrawlProxy(url + htmlquery.InnerText(page))
				}
			}
		})

		// c.OnHTML()
		c.OnScraped(func(r *colly.Response) {
			fmt.Println("Finished:", r.Request.URL)
		})
		c.Visit(url)
		c.Wait()
		p.PushProxy(proxyList...)
		<-timeTickerChan
	}
}
