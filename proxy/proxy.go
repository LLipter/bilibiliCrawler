package proxy

import (
	"encoding/json"
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
func GetProxies() []string {
	apiAddr := "https://dev.kdlapi.com/api/getproxy/?orderid=904212196080767&num=1000&b_pcchrome=1&b_pcie=1&b_pcff=1&protocol=2&method=2&an_tr=1&an_an=1&an_ha=1&sp1=1&sp2=1&quality=1&format=json&sep=1"

	for {
		resp, err := http.Get(apiAddr)
		if err != nil {
			log.Fatalln("get proxy failed, " + err.Error())
			// retry after 10 seconds
			time.Sleep(time.Second * 10)
		}

		var proxyObj proxyJson
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("read failed, " + err.Error())
			// retry after 10 seconds
			time.Sleep(time.Second * 10)
			continue
		}
		resp.Body.Close()

		err = json.Unmarshal(buf, &proxyObj)
		if err != nil {
			log.Fatalln(err)
			// retry after 10 seconds
			time.Sleep(time.Second * 10)
			continue
		}

		if proxyObj.Code != 0 {
			log.Fatalln("illegal get proxy parameters")
			// retry after 10 seconds
			time.Sleep(time.Second * 10)
			continue
		}

		ProxyPool = proxyObj.Data.Proxy_list

		// refresh proxy pool every minute
		time.Sleep(time.Minute)
	}
}
