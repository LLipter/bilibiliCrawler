package proxy

import (
	"encoding/json"
	"errors"
	"github.com/LLipter/bilibiliCrawler/conf"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	ProxyPool []string
)

type data struct {
	Count      int
	Proxy_list []string
}

type proxyJson struct {
	Msg  string
	Code int
	Data data
}

// change to your own codes to get proxy
func GetProxy() error {
	apiAddr := "http://dev.kdlapi.com/api/getproxy/?orderid=904212196080767&num=1000&b_pcchrome=1&b_pcie=1&b_pcff=1&protocol=1&method=1&an_ha=1&sp1=1&sp2=1&quality=1&format=json&sep=1"
	resp, err := http.Get(apiAddr)
	if err != nil {
		return errors.New("get proxy failed, " + err.Error())
	}

	var proxyObj proxyJson
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	err = json.Unmarshal(buf, &proxyObj)
	if err != nil {
		return err
	}

	if proxyObj.Code != 0 {
		return errors.New("illegal get proxy parameters")
	}

	ProxyPool = proxyObj.Data.Proxy_list
	return nil
}

func GetProxyRoutine() {
	for t := 0; t < conf.NetworkConfig.RetryTimes; t++ {
		err := GetProxy()
		if err != nil {
			log.Println(err)
			continue
		}
		t = 0

		// refresh proxy pool every 30 second
		time.Sleep(time.Second * 30)
	}
	log.Fatal("cannot get proxy")
}
