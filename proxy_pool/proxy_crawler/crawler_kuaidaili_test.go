package proxy_crawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestCrawlProxy(t *testing.T) {
	c := &CrawlerKuaidaili{}
	for _, ip := range c.CrawlProxy() {
		fmt.Println(CheckProxy(ip))
	}
}

func CheckProxy(proxyAddr string) bool {
	httpUrl := "http://icanhazip.com"
	proxy, _ := url.Parse(proxyAddr)

	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}
	httpClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	res, err := httpClient.Get(httpUrl)
	if err != nil {
		// fmt.Println("错误信息：",err)
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Println(err)
		return false
	}
	c, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK || string(c) == "" {
		return false
	}
	return true
}
