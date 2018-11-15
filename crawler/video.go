package crawler

import (
	"encoding/json"
	"github.com/LLipter/bilibiliCrawler/conf"
	"github.com/LLipter/bilibiliCrawler/util"
	"github.com/LLipter/bilibiliCrawler/util/db"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func CrawlVideo(startAid int, endAid int) {
	for i := startAid; i <= endAid; i++ {
		// control max number of crawler go routine
		curCrawlerNo <- true
		go videoCrawlerRoutine(i)
	}

	// make sure all crawler go routine finish their job
	wg.Wait()

}

func videoCrawlerRoutine(aid int) {
	wg.Add(1)
	defer wg.Done()
	defer func() {
		<-curCrawlerNo
	}()

	var err error
	for t := 0; t < conf.NetworkConfig.RetryTimes; t++ {
		err = getVideoData(aid)
		if err == nil {
			return
		}
	}

	// failed with unknown reason
	if err != nil {
		log.Printf("aid=%d crawler failed, %v\n", aid, err)
	}
	var video conf.Video
	video.Status = 2
	video.Aid = int64(aid)
	err = db.InsertVideo(video)
	if err != nil {
		log.Printf("aid=%d insertion failed, %v\n", aid, err)
	}
}

func getVideoData(aid int) error {
	var data conf.Info
	err := getVideoBasicData(aid, &data)
	if err != nil {
		return err
	}

	// if this video exists, crawl more data
	if data.Code == 0 {
		err = getVideoMoreData(aid, &data.Data)
		if err != nil {
			return err
		}
	} else {
		data.Data.Status = 1
		data.Data.Aid = int64(aid)
	}

	err = db.InsertVideo(data.Data)
	if err != nil {
		return err
	}

	return nil
}

func getVideoBasicData(aid int, data *conf.Info) error {
	addr := "http://api.bilibili.com/archive_stat/stat?aid=" + strconv.Itoa(aid)
	resp, err := getResp(addr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, data)
	if err != nil {
		if strings.HasPrefix(err.Error(), "json: cannot unmarshal string into Go struct field") {
			data.Code = -1
		} else {
			return err
		}
	}

	return nil
}

func getVideoMoreData(aid int, video *conf.Video) error {
	addr := "http://api.bilibili.com/view?appkey=8e9fc618fbd41e28&id=" + strconv.Itoa(aid)
	resp, err := getResp(addr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var jsonObj map[string]interface{}
	json.Unmarshal(buf, &jsonObj)

	// get tid
	video.Tid, err = util.JsonGetInt64(jsonObj, "tid")
	if err != nil {
		return err
	}

	// get title
	video.Title, err = util.JsonGetStr(jsonObj, "title")
	if err != nil {
		return err
	}

	// get mid
	video.Mid, err = util.JsonGetInt64(jsonObj, "mid")
	if err != nil {
		return err
	}

	// get pubdate
	pubdate, err := util.JsonGetInt64(jsonObj, "created")
	if err != nil {
		return err
	}
	video.Pubdate = time.Unix(pubdate, 0)

	// get cid
	video.Cid, err = util.JsonGetInt64(jsonObj, "cid")
	if err != nil {
		return err
	}

	return nil
}
