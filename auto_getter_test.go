package main

import "testing"

func TestDecorator(t *testing.T) {
	var c AutoProxyGetter = &OriAutoProxyGetter{}
	c = WrapWithTimeDecorator(c, 100)
	c.CrawlProxy("http://www.ip3366.net/?stype=1&page=1")
}
