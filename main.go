package main

func main() {
	var c AutoProxyGetter = &OriAutoProxyGetter{}
	c = WrapWithTimeDecorator(c, 180)
	c = WrapWithThresholdDecorator(c, 80)
	c.CrawlProxy("http://www.ip3366.net/?stype=1&page=1")
}
