package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/LLipter/bilibili-report/util"
)


func sendRequest(addr string, useProxy bool) (util.Info, error) {
	var resp *http.Response
	var err error
	if useProxy{
		// TODO: add proxies pool
		urlproxy, err := url.Parse("http://183.245.99.52:80")
		if err != nil {
			return util.Info{}, err
		}
		client := http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlproxy),
			},
		}
		resp, err = client.Get(addr)
		if err != nil {
			return util.Info{}, err
		}
	}else{
		resp, err = http.Get(addr)
		if err != nil {
			return util.Info{}, err
		}
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return util.Info{}, err
	}

	var info util.Info
	err = json.Unmarshal(data, &info)
	if err != nil {
		return util.Info{}, err
	}

	return info, nil
}

func init(){
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("init successfully. ")
	log.Println("Other logs will be store in log file. No information will appear here")
	//log.SetOutput()
}

func spider(aid int) error{
	info, err := sendRequest("https://api.bilibili.com/archive_stat/stat?aid=" + strconv.Itoa(aid), false)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", info)
	var video util.Video

	if info.Code != 0 {
		video.Status = 1
		video.Aid = aid
	} else {
		video = info.Data
	}

	err = util.InsertVideo(video)
	if err != nil{
		return err
	}

	return nil
}

func main() {

	for i:=1;i<=10;i++{
		err := spider(i)
		if err != nil{
			fmt.Println(err)
		}
	}

	util.CloseDatabase()

}
