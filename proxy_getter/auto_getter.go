package proxy_pool

import (
	"time"

	"github.com/Apale7/lazy_proxy/proxy_getter/proxy_crawler"
)

type AutoProxyGetter interface {
	ProxyPool
	proxy_crawler.Crawler
}

type DefaultAutoProxyGetter struct {
	ProxyPool
	proxy_crawler.Crawler
}

type WithTimeDecorator struct {
	AutoProxyGetter
	interval int
}

// WrapWithTimeDecorator 周期性爬取代理ip
func WrapWithTimeDecorator(a AutoProxyGetter, interval int) *WithTimeDecorator {
	getter := &WithTimeDecorator{
		AutoProxyGetter: a,
		interval:        interval,
	}
	go func() {
		timeTickerChan := time.NewTicker(time.Second * time.Duration(getter.interval))
		for {
			proxyList := getter.AutoProxyGetter.CrawlProxy()
			getter.AutoProxyGetter.PushProxy(proxyList...)
			<-timeTickerChan.C
		}
	}()
	return getter
}

type WithThresholdDecorator struct {
	AutoProxyGetter
	threshold int
}

// WrapWithThresholdDecorator 达到阈值时爬取代理ip
func WrapWithThresholdDecorator(a AutoProxyGetter, threshold int) *WithThresholdDecorator {
	getter := &WithThresholdDecorator{
		AutoProxyGetter: a,
		threshold:       threshold,
	}
	go func() {
		for {
			if getter.AutoProxyGetter.LenOfProxies() < getter.threshold {
				proxyList := getter.AutoProxyGetter.CrawlProxy()
				getter.AutoProxyGetter.PushProxy(proxyList...)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	return getter
}
