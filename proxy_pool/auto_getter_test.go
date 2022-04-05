package proxy_pool

import (
	"fmt"
	"testing"
	"time"

	"github.com/Apale7/lazy_proxy/proxy_pool/proxy_crawler"
)

func TestDecorator(t *testing.T) {
	var c AutoProxyGetter = &DefaultAutoProxyGetter{
		ProxyPool: &DefaultProxyPool{},
		Crawler:   &proxy_crawler.CrawlerKuaidaili{},
	}
	c = WrapWithTimeDecorator(c, 300)
	c = WrapWithThresholdDecorator(c, 25)

	for c.LenOfProxies() == 0 {
		time.Sleep(time.Second * 2)
	}
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
