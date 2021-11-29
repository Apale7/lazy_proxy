package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
)

type AutoProxyGetter interface {
	ProxyGetter
	CrawlProxy(url string) // Crawl available proxies from the Internet
}

type DefaultAutoProxyGetter struct {
	DefaultProxyGetter
	interval int // How many seconds between crawling proxies
}

// 先加url,再爬proxy

func (p *DefaultAutoProxyGetter) CrawlProxy(url string) {
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
			// nextPage
			// pages := htmlquery.Find(doc, `//*[@id="listnav"]//a/@href`)
			// for _, page := range pages {
			// 	fmt.Println("Page:", htmlquery.InnerText(page))
			// }
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
