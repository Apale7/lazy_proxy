package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly/v2"
)

type Crawler interface {
	CrawlProxy() []string // Crawl available proxies from the Internet
}

type CrawlerIP3366 struct {
	url string
}

func (cl *CrawlerIP3366) CrawlProxy() []string {
	// 爬最大页码
	//.......

	// 遍历每一页
	// go func(){爬一页}
	var i int64 = 2
	// Init collyCollector
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
	// flower
	// if d.proxies != nil && len(d.proxies) != 0 {
	// 	ipPort := []string{}
	// 	for _, v := range d.proxies {
	// 		ipPort = append(ipPort, "socks5://"+v)
	// 	}
	// 	rp, err := proxy.RoundRobinProxySwitcher(ipPorgetter...)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	c.SetProxyFunc(rp)
	// }
	proxyList := []string{}
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Fatal(err)
	})
	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader((string(r.Body))))
		if err != nil {
			log.Fatal(err)
		}
		nodes := htmlquery.Find(doc, `//tbody/tr`)
		for _, node := range nodes {
			addr := htmlquery.FindOne(node, `./td[1]/text()`)
			port := htmlquery.FindOne(node, `./td[2]/text()`)
			proxy := htmlquery.InnerText(addr) + ":" + htmlquery.InnerText(port)
			proxyList = append(proxyList, proxy)
		}
	})
	c.OnHTML("#listnav > ul > b > font", func(e *colly.HTMLElement) {
		max, _ := strconv.ParseInt(e.Text, 10, 64)
		max = max / 10
		for i <= max {
			num := strconv.FormatInt(i, 10)
			e.Request.Visit("http://www.ip3366.net/?stype=1&page=" + num)
			i++
		}
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished:", r.Request.URL)
	})
	c.Visit(cl.url)
	c.Wait()
	return proxyList
}
type AutoProxyGetter interface {
	ProxyGetter
	Crawler

}
type DefaultAutoProxyGetter struct {
	crawler Crawler
	ProxyGetter
}

type WithTimeDecorator struct {
	AutoProxyGetter
	interval int
}

func WrapWithTimeDecorator(a AutoProxyGetter, interval int) *WithTimeDecorator {
	getter := &WithTimeDecorator{
		AutoProxyGetter: a,
		interval:        interval,
	}
	go func() {
		timeTickerChan := time.Tick(time.Second * time.Duration(getter.interval))
		for {
			getter.AutoProxyGetter.CrawlProxy()
			<-timeTickerChan
		}
	}()
	return getter
}

type WithThresholdDecorator struct {
	AutoProxyGetter
	threshold int
}

func WrapWithThresholdDecorator(a AutoProxyGetter, threshold int) *WithThresholdDecorator {
	getter := &WithThresholdDecorator{
		AutoProxyGetter: a,
		threshold:       threshold,
	}
	go func() {
		for getter.AutoProxyGetter.LenOfProxies() < getter.threshold {
			time.Sleep(1 * time.Minute)
			getter.AutoProxyGetter.crawler.CrawlProxy()
		}
	}()
	return getter
}

func (g *WithThresholdDecorator) EraseProxy(proxy string) int {
	num := g.AutoProxyGetter.EraseProxy(proxy)
	if num < g.threshold {
		go g.AutoProxyGetter.crawler.CrawlProxy()
	}
	return num
}
