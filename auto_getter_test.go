package main

import "testing"

func TestAutoCrawl(t *testing.T) {
	p := &TimerAutoProxyGetter{
		interval: 10,
	}
	p.CrawlProxy("http://www.ip3366.net/")
}
