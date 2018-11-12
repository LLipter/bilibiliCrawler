package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/LLipter/bilibili-report/util"
)

type Video struct {
	Aid        int
	Status     int
	View       int
	Dannmaku   int
	Reply      int
	Favorite   int
	Coin       int
	Share      int
	Now_rank   int
	His_rank   int
	Support    int
	Dislike    int
	No_reprint int
	Copyright  int
}

type Info struct {
	Code    int
	Message string
	Data    Video
}

func sendRequest(addr string) (Info, error) {
	//urlproxy, err := url.Parse("http://183.245.99.52:80")
	//if err != nil {
	//	return err
	//}
	//
	//client := http.Client{
	//	Transport: &http.Transport{
	//		Proxy: http.ProxyURL(urlproxy),
	//	},
	//}
	//resp, err := client.Get(addr)
	//if err != nil {
	//	return err
	//}

	resp, err := http.Get(addr)
	if err != nil {
		return Info{}, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Info{}, err
	}

	var info Info
	err = json.Unmarshal(data, &info)
	if err != nil {
		return Info{}, err
	}

	return info,nil
}



func main() {
	info,err := sendRequest("https://api.bilibili.com/archive_stat/stat?aid=2")
	if err != nil {
		fmt.Println(err)
	}

	util.PrintJson(info)
	fmt.Println(info)

}
