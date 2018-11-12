package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/LLipter/bilibili-report/util"
)

var(
	retryTimes = 3
	wg sync.WaitGroup
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

func crawler(aid int) error{
	info, err := sendRequest("https://api.bilibili.com/archive_stat/stat?aid=" + strconv.Itoa(aid), false)
	if err != nil {
		return err
	}

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

func crawlerRoutine(aid int){
	defer wg.Done()
	for t:=0;t<retryTimes;t++{
		err := crawler(aid)
		if err == nil{
			return
		}else{
			log.Printf("aid=%d crawler failed, %v\n", aid, err)
		}
	}
	// failed with unknown reason
	var video util.Video
	video.Status = 2
	video.Aid = aid
	err := util.InsertVideo(video)
	if err != nil{
		log.Printf("aid=%d insertion failed, %v\n", aid, err)
	}
}

func main() {

	for i:=1;i<=100;i++{
		wg.Add(1)
		go crawlerRoutine(i)
	}

	wg.Wait()

	util.CloseDatabase()

}
