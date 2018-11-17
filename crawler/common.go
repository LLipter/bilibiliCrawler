package crawler

import (
	"errors"
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
	curCrawlerNo = make(chan bool, conf.NetworkConfig.MaxCrawlerNum)
}

func getResp(addr string) ([]byte, error) {
	tr := &http.Transport{
		// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#close_http_conn
		/*
			Close connection immediately,
			otherwise the number of opened file will keep growing.
		*/
		DisableKeepAlives: true,
	}
	if conf.NetworkConfig.UseProxy {
		length := len(proxy.ProxyPool)
		if length == 0 {
			return nil, errors.New("no proxy")
		}
		proxyAddr := proxy.ProxyPool[rand.Intn(length)]
		urlProxy, err := url.Parse("http://" + proxyAddr)
		if err != nil {
			return nil, err
		}
		tr.Proxy = http.ProxyURL(urlProxy)
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 20,
	}

	req, err := http.NewRequest("GET", addr, nil)
	req.Header.Add("User-Agent", conf.NetworkConfig.UserAgent)
	resp, err := client.Do(req)
	// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#close_http_resp_body
	/*
		Most of the time when your http request fails the `resp` variable will be nil
		and the `err` variable will be non-nil.
		However, when you get a redirection failure both variables will be non-nil.
		This means you can still end up with a leak.
	*/
	if resp != nil {
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
