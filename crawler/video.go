package crawler

import (
	"encoding/json"
	"github.com/LLipter/bilibiliVideoDataCrawler/conf"
	"github.com/LLipter/bilibiliVideoDataCrawler/util"
	"github.com/LLipter/bilibiliVideoDataCrawler/util/db"
	"github.com/PuerkitoBio/goquery"
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
		<- curCrawlerNo
	}()

	var err error
	for t := 0; t < conf.RetryTimes; t++ {
		err = getVideoData(aid)
		if err == nil {
			return
		}
	}

	if err != nil {
		log.Printf("aid=%d crawler failed, %v\n", aid, err)
	}

	// failed with unknown reason
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
	}

	var video conf.Video
	if data.Code != 0 {
		video.Status = 1
		video.Aid = int64(aid)
	} else {
		video = data.Data
	}

	err = db.InsertVideo(video)
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
	addr := "http://www.bilibili.com/video/av" + strconv.Itoa(aid)
	resp, err := getResp(addr)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	var jsonstr string
	doc.Find("script").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		content := selection.Text()
		if strings.HasPrefix(content, "window.__INITIAL_STATE__=") {
			jsonstr = content[25 : len(content)-122]
			return false
		}
		return true
	})

	var jsonObj map[string]interface{}
	json.Unmarshal([]byte(jsonstr), &jsonObj)

	videoJson, err := util.JsonGetDict(jsonObj, "videoData")
	if err != nil {
		return err
	}

	// get title
	video.Title, err = util.JsonGetStr(videoJson, "title")
	if err != nil {
		return err
	}

	// get pubdate
	pubdate, err := util.JsonGetInt64(videoJson, "pubdate")
	if err != nil {
		return err
	}
	video.Pubdate = time.Unix(pubdate, 0)

	// get duration
	video.Duration, err = util.JsonGetInt64(videoJson, "duration")
	if err != nil {
		return err
	}

	// get ownerJson id
	ownerJson, err := util.JsonGetDict(videoJson, "owner")
	if err != nil {
		return err
	}
	video.Owner, err = util.JsonGetInt64(ownerJson, "mid")
	if err != nil {
		return err
	}

	// get pages
	pages, err := util.JsonGetArray(videoJson, "pages")
	if err != nil {
		return err
	}
	for _, pageObj := range pages {
		pageJson, ok := pageObj.(map[string]interface{})
		if !ok {
			return util.TypeError("pages")
		}

		// get chatid
		var page conf.Page
		page.Chatid, err = util.JsonGetInt64(pageJson, "cid")
		if err != nil {
			return err
		}

		// get duration
		page.Duration, err = util.JsonGetInt64(pageJson, "duration")
		if err != nil {
			return err
		}

		// get pageNo
		page.PageNo, err = util.JsonGetInt64(pageJson, "page")
		if err != nil {
			return err
		}

		// get subtititle
		page.Subtitle, err = util.JsonGetStr(pageJson, "part")
		if err != nil {
			return err
		}

		video.Pages = append(video.Pages, page)

	}

	return nil
}
