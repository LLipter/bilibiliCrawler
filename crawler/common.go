package crawler

import (
	"github.com/LLipter/bilibili-report/conf"
	"github.com/LLipter/bilibili-report/proxy"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

var (
	wg sync.WaitGroup
)

func getResp(addr string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", addr, nil)
	req.Header.Add("User-Agent", conf.UserAgent)
	if conf.UseProxy {
		length := len(proxy.ProxyPool)
		proxyAddr := proxy.ProxyPool[rand.Intn(length)]
		urlProxy, err := url.Parse("http://" + proxyAddr)
		if err != nil {
			return nil, err
		}
		client = http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}