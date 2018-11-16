package crawler

import (
	"encoding/json"
	"errors"
	"github.com/LLipter/bilibiliCrawler/conf"
	"github.com/LLipter/bilibiliCrawler/util/db"
	"time"
)

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
