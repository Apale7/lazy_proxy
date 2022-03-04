package proxy_crawler

import (
	"fmt"
	"testing"
)

func TestCrawlProxy(t *testing.T) {
	c := &CrawlerKuaidaili{}
	for i, ip := range c.CrawlProxy() {
		fmt.Printf("%d: %s\n", i, ip)
	}
}
