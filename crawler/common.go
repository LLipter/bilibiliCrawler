package crawler

import (
	"github.com/LLipter/bilibiliCrawler/conf"
	"github.com/LLipter/bilibiliCrawler/proxy"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	wg           sync.WaitGroup
	curCrawlerNo chan bool
)

func init() {
	curCrawlerNo = make(chan bool, conf.VideoCrawlerConfig.MaxCrawlerNum)
}

func getResp(addr string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", addr, nil)
	req.Header.Add("User-Agent", conf.NetworkConfig.UserAgent)
	if conf.NetworkConfig.UseProxy {
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
			Timeout: time.Second * 20,
		}
	}
	resp, err := client.Do(req)
	// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#close_http_resp_body
	/*
	Most of the time when your http request fails the `resp` variable will be nil
	and the `err` variable will be non-nil.
	However, when you get a redirection failure both variables will be non-nil.
	This means you can still end up with a leak.
	 */
	if resp != nil{
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}


	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
