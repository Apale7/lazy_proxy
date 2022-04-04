package proxy_crawler

type Crawler interface {
	CrawlProxy() []string // Crawl available proxies from the Internet
}
