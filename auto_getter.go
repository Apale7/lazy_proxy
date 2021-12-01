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

type AutoProxyGetter interface {
	ProxyGetter
	CrawlProxy(url string) // Crawl available proxies from the Internet
}

type OriAutoProxyGetter struct {
	DefaultProxyGetter
}

func (d *OriAutoProxyGetter) CrawlProxy(url string) {
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
	// 	rp, err := proxy.RoundRobinProxySwitcher(ipPort...)
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
	c.Visit(url)
	c.Wait()
	d.PushProxy(proxyList...)
	fmt.Println("LENGTH OF PROXIES:", len(d.proxies))
}

type WithTimeDecorator struct {
	AutoProxyGetter
	interval int
}

func WrapWithTimeDecorator(a AutoProxyGetter, interval int) AutoProxyGetter {
	return &WithTimeDecorator{
		AutoProxyGetter: a,
		interval:        interval,
	}
}

func (t *WithTimeDecorator) CrawlProxy(url string) {
	timeTickerChan := time.Tick(time.Second * time.Duration(t.interval))
	for {
		t.AutoProxyGetter.CrawlProxy(url)
		<-timeTickerChan
	}
}

type WithThresholdDecorator struct {
	AutoProxyGetter
	threshold int
}

func WrapWithThresholdDecorator(a AutoProxyGetter, threshold int) AutoProxyGetter {
	return &WithThresholdDecorator{
		AutoProxyGetter: a,
		threshold:       threshold,
	}
}

func (t *WithThresholdDecorator) CrawlProxy(url string) {
	for t.AutoProxyGetter.LenOfProxies() < t.threshold {
		t.AutoProxyGetter.CrawlProxy(url)
	}
}
