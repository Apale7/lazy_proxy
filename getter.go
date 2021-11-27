package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type ProxyGetter interface {
	GetProxy(string, error)       // return an usable proxy. if there is not an usable proxy, return "", error
	CheckProxy(proxy string) bool // check if the proxy is usable
	EraseProxy(proxy string)      // erase the proxy from the proxy list
	PushProxy(proxy string)       // push the proxy into the proxy list
}

type DefaultProxyGetter struct {
	now                 int // return an proxy in sequence, without using random numbers
	proxys              []string
	CheckBeforeGetProxy bool // if true, proxy will be checked for availability before returning; otherwise, it will be returned directly
	lock                sync.Mutex
}

func (p *DefaultProxyGetter) GetProxy() (string, error) {
	if p.proxys == nil || len(p.proxys) == 0 {
		return "", nil
	}
	proxy := p.proxys[p.now]
	p.now = (p.now + 1) % len(p.proxys)
	if p.CheckBeforeGetProxy {
		for len(p.proxys) > 0 && !p.CheckProxy(proxy) {
			p.EraseProxy(proxy)
			proxy = p.proxys[p.now]
			p.now = (p.now + 1) % len(p.proxys)
		}
	}
	if p.proxys == nil || len(p.proxys) == 0 {
		return "", nil
	}
	return proxy, nil
}

// The efficiency is not high when the number of proxies is large
func (p *DefaultProxyGetter) EraseProxy(proxy string) {
	p.lock.Lock()
	for i, v := range p.proxys {
		if v == proxy {
			p.proxys = append(p.proxys[:i], p.proxys[i+1:]...)
			break
		}
	}
	p.lock.Unlock()
}

func (p *DefaultProxyGetter) PushProxy(proxy string) {
	p.lock.Lock()
	if p.CheckProxy(proxy) {
		p.proxys = append(p.proxys, proxy)
	}
	p.lock.Unlock()
}

func (p *DefaultProxyGetter) CheckProxy(proxyAddr string) bool {
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
		//fmt.Println("错误信息：",err)
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
