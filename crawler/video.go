package crawler

import (
	"encoding/json"
	"github.com/LLipter/bilibili-report/conf"
	"github.com/LLipter/bilibili-report/util"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func CrawlVideo(startAid int, endAid int) {
	for i := startAid; i <= endAid; i++ {
		for runtime.NumGoroutine() > conf.MaxGoroutinueNum {
			time.Sleep(time.Second)
		}
		wg.Add(1)
		go videoCrawlerRoutine(i)
	}

	wg.Wait()
}

func videoCrawlerRoutine(aid int) {
	defer wg.Done()
	for t := 0; t < conf.RetryTimes; t++ {
		err := getVideoData(aid)
		if err == nil {
			return
		} else {
			log.Printf("aid=%d crawler failed, %v\n", aid, err)
		}
	}
	// failed with unknown reason
	var video util.Video
	video.Status = 2
	video.Aid = int64(aid)
	err := util.InsertVideo(video)
	if err != nil {
		log.Printf("aid=%d insertion failed, %v\n", aid, err)
	}
}

func getVideoData(aid int) error {
	var data util.Info
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

	var video util.Video
	if data.Code != 0 {
		video.Status = 1
		video.Aid = int64(aid)
	} else {
		video = data.Data
	}

	err = util.InsertVideo(video)
	if err != nil {
		return err
	}

	return nil
}

func getVideoBasicData(aid int, data *util.Info) error {
	addr := "https://api.bilibili.com/archive_stat/stat?aid=" + strconv.Itoa(aid)
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

func getVideoMoreData(aid int, video *util.Video) error {
	addr := "https://www.bilibili.com/video/av" + strconv.Itoa(aid)
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
		var page util.Page
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
