package proxy_crawler

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

func (craw *CrawlerKuaidaili) newCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36"),
		colly.AllowURLRevisit(),
	)

	var proxy func(*http.Request) (*url.URL, error)
	if len(craw.proxys) > 0 {
		p, _ := url.Parse(craw.proxys[rand.Intn(len(craw.proxys))])
		proxy = http.ProxyURL(p)
	}

	c.WithTransport(&http.Transport{
		Proxy: proxy,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 30 * time.Second,
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Fatal(err)
	})
	return c
}

const (
	maxPage int64 = 100
)

type CrawlerKuaidaili struct {
	mutex  sync.Mutex
	proxys []string
}

func (craw *CrawlerKuaidaili) CrawlProxy() (proxy []string) {
	if !craw.mutex.TryLock() {
		return
	}
	defer craw.mutex.Unlock()
	c := craw.newCollector()
	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader((string(r.Body))))
		if err != nil {
			log.Fatal(err)
		}
		portNodes := htmlquery.Find(doc, `//td[@data-title="PORT"]`)
		IPNodes := htmlquery.Find(doc, `//td[@data-title="IP"]`)
		proxy = make([]string, len(IPNodes))
		for i, node := range IPNodes {
			proxy[i] = htmlquery.InnerText(node) + ":" + htmlquery.InnerText(portNodes[i])
		}
	})

	if err := c.Visit(fmt.Sprintf("https://www.kuaidaili.com/free/inha/%d/", rand.Int31n(int32(maxPage)))); err != nil {
		logrus.Error(err)
	}
	c.Wait()
	craw.proxys = proxy

	return proxy
}
