package crawler

import (
	"github.com/LLipter/bilibili-report/conf"
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
		// TODO: add proxies pool
		urlproxy, err := url.Parse("http://183.245.99.52:80")
		if err != nil {
			return nil, err
		}
		client = http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlproxy),
			},
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
