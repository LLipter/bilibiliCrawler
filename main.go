package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/LLipter/bilibili-report/util"
)


func sendRequest(addr string) (util.Info, error) {
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
		return util.Info{}, err
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

func main() {
	info, err := sendRequest("https://api.bilibili.com/archive_stat/stat?aid=12")
	if err != nil {
		fmt.Println(err)
	}

	util.PrintJson(info)
	fmt.Printf("%+v\n", info)

	if info.Code != 0 {
		fmt.Println(info.Message)
	} else {
		err := util.InsertVideo(info.Data)
		if err != nil{
			fmt.Println(err)
		}
	}

	util.CloseDatabase()

}
