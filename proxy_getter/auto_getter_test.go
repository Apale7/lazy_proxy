package proxy_getter

import (
	"fmt"
	"testing"
	"time"

	"github.com/Apale7/lazy_proxy/proxy_getter/proxy_crawler"
)

func TestDecorator(t *testing.T) {
	var c AutoProxyGetter = &DefaultAutoProxyGetter{
		ProxyGetter: &DefaultProxyGetter{},
		Crawler:     &proxy_crawler.CrawlerIP3366{},
	}
	c = WrapWithTimeDecorator(c, 300)
	c = WrapWithThresholdDecorator(c, 80)
	proxyList := c.CrawlProxy()
	c.PushProxy(proxyList...)
	go func() {
		for {
			p, _ := c.GetProxy()
			fmt.Printf("c.EraseProxy(p): %v\n", c.EraseProxy(p))
			time.Sleep(time.Millisecond * 5000)
		}
	}()
	select {}
}

func TestCrawlerIP3366(t *testing.T) {
	c := &proxy_crawler.CrawlerIP3366{}
	c.CrawlProxy()
}
