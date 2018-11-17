package crawler

import (
	"encoding/json"
	"errors"
	"github.com/LLipter/bilibiliCrawler/conf"
	"github.com/LLipter/bilibiliCrawler/util/db"
	"log"
	"time"
)

func CrawOnline() {
	for t := 0; t < conf.NetworkConfig.RetryTimes; t++ {
		err := getOnlineData()
		if err != nil {
			log.Println(err)
			continue
		}
		t = 0

		// crawl data every second
		time.Sleep(time.Minute)
	}
	log.Fatalln("cannot get online data")
}

func getOnlineData() error {
	addr := "http://api.bilibili.com/x/web-interface/online"
	buf, err := getResp(addr)
	if err != nil {
		return err
	}

	var onlineJson conf.OnlineJson
	err = json.Unmarshal(buf, &onlineJson)
	if err != nil {
		return err
	}

	if onlineJson.Code != 0 {
		return errors.New(onlineJson.Message)
	}

	onlineJson.Timestamp = time.Now()
	err = db.InsertOnline(onlineJson)
	if err != nil {
		return err
	}
	return nil
}
