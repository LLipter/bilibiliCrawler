package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LLipter/bilibiliCrawler/conf"
	"github.com/LLipter/bilibiliCrawler/util"
	"github.com/LLipter/bilibiliCrawler/util/db"
	"log"
	"strconv"
	"time"
)

func CrawlBangumi() {
	var err error
	var totalNo int
	for t := 0; t < conf.NetworkConfig.RetryTimes; t++ {
		totalNo, err = getTotalNumber()
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Println("cannot get total number: " + err.Error())
	}

	pageSize := 20
	pageNo := totalNo / pageSize
	if totalNo%pageSize != 0 {
		pageNo++
	}

	for page := 1; page <= pageNo; page++ {
		// control max number of crawler go routine
		curCrawlerNo <- true
		wg.Add(1)
		go pageCrawlerRoutine(page)
	}

	wg.Wait()
}

func getTotalNumber() (int, error) {
	addr := "https://bangumi.bilibili.com/media/web_api/search/result?season_version=1&area=2&is_finish=1&copyright=-1&season_status=-1&season_month=-1&pub_date=-1&style_id=-1&order=3&st=1&sort=0&page=1&season_type=1&pagesize=20"
	buf, err := getResp(addr)
	if err != nil {
		return 0, err
	}

	var bangumiJson conf.BangumiJson
	err = json.Unmarshal(buf, &bangumiJson)
	if err != nil {
		return 0, err
	}

	return int(bangumiJson.Data.Page.Total), nil
}

func pageCrawlerRoutine(page int) {
	defer wg.Done()
	defer func() {
		<-curCrawlerNo
	}()

	var err error
	for t := 0; t < conf.NetworkConfig.RetryTimes; t++ {
		err = getPage(page)
		if err == nil {
			return
		}
	}

	log.Println(fmt.Sprintf("cannot get page %d: ", page) + err.Error())

}

func getPage(page int) error {
	addr := "https://bangumi.bilibili.com/media/web_api/search/result?season_version=1&area=2&is_finish=1&copyright=-1&season_status=-1&season_month=-1&pub_date=-1&style_id=-1&order=3&st=1&sort=0&season_type=1&pagesize=20&page=" + strconv.Itoa(page)
	buf, err := getResp(addr)
	if err != nil {
		return err
	}

	var bangumiJson conf.BangumiJson
	err = json.Unmarshal(buf, &bangumiJson)
	if err != nil {
		return err
	}

	for idx, data := range bangumiJson.Data.Bangumis {
		bangumi, err := convertBangumi(data)
		if err != nil {
			log.Println(fmt.Sprintf("page=%d,idx=%d bangumi crawler failed: ", page, idx) + err.Error())
			continue
		}

		// control max number of crawler go routine
		fmt.Printf("start crawl: page=%d, index=%d, bangumi=%s\n", page, idx, bangumi.Title)
		curCrawlerNo <- true
		wg.Add(1)
		go eplistCrawlerRoutine(bangumi)
	}

	return nil
}

func eplistCrawlerRoutine(bangumi conf.Bangumi) {
	defer wg.Done()
	defer func() {
		<-curCrawlerNo
	}()

	var err error
	for t := 0; t < conf.NetworkConfig.RetryTimes; t++ {
		// get eplist
		err = getEplist(int(bangumi.Sid), bangumi)
		if err == nil {
			return
		}
	}

	// failed with unknown reason
	if err != nil {
		log.Printf("bangumi=%s crawler failed, %v\n", bangumi.Title, err)
	}
}

func getEplist(sid int, bangumi conf.Bangumi) error {
	addr := "http://bangumi.bilibili.com/web_api/get_ep_list?season_type=1?&season_id=" + strconv.Itoa(sid)
	buf, err := getResp(addr)
	if err != nil {
		return errors.New(fmt.Sprintf("sid=%d cannot get eplist: ", sid) + err.Error())
	}

	var eplistJson conf.EplistJson
	err = json.Unmarshal(buf, &eplistJson)
	if err != nil {
		return errors.New(fmt.Sprintf("sid=%d cannot get eplist: ", sid) + err.Error())
	}

	for i := 0; i < len(eplistJson.Data); i++ {
		aid := int(eplistJson.Data[i].Aid)
		var videoJson conf.VideoJson
		for t := 0; t < conf.NetworkConfig.RetryTimes; t++ {
			err = getVideoBasicData(aid, &videoJson)
			if err == nil {
				break
			}
		}
		if err != nil {
			return errors.New(fmt.Sprintf("aid=%d cannot get view: ", aid) + err.Error())
		}
		eplistJson.Data[i].View = videoJson.Data.View
	}

	bangumi.Eplist = eplistJson.Data

	err = db.InsertBangumi(bangumi)
	if err != nil {
		return err
	}

	return nil
}

func convertBangumi(data interface{}) (conf.Bangumi, error) {
	dataDict, ok := data.(map[string]interface{})
	if !ok {
		return conf.Bangumi{}, errors.New("bangumi type error")
	}
	var err error
	var bangumi conf.Bangumi
	// get media_id
	bangumi.MediaID, err = util.JsonGetInt64(dataDict, "media_id")
	if err != nil {
		return conf.Bangumi{}, err
	}

	// get season id
	bangumi.Sid, err = util.JsonGetInt64(dataDict, "season_id")
	if err != nil {
		return conf.Bangumi{}, err
	}

	// get season title
	bangumi.Title, err = util.JsonGetStr(dataDict, "title")
	if err != nil {
		return conf.Bangumi{}, err
	}

	dataDict, err = util.JsonGetDict(dataDict, "order")
	if err != nil {
		return conf.Bangumi{}, err
	}

	// get pubdate
	pubdate, err := util.JsonGetInt64(dataDict, "pub_date")
	if err != nil {
		return conf.Bangumi{}, err
	}
	bangumi.Pubdate = time.Unix(pubdate, 0)

	// get follow
	followStr, err := util.JsonGetStr(dataDict, "follow")
	if err != nil {
		return conf.Bangumi{}, err
	}
	var baseNumber float64
	var suffix string
	n, err := fmt.Sscanf(followStr, "%f%s", &baseNumber, &suffix)
	if err != nil {
		return conf.Bangumi{}, err
	}
	if n != 2 {
		return conf.Bangumi{}, errors.New(followStr + " type error")
	}
	if suffix == "万人追番" {
		bangumi.Follow = int64(baseNumber * 10000)
	} else if suffix == "人追番" {
		bangumi.Follow = int64(baseNumber)
	} else {
		return conf.Bangumi{}, errors.New("unknown suffix: " + suffix)
	}

	// get view
	viewStr, err := util.JsonGetStr(dataDict, "play")
	if err != nil {
		return conf.Bangumi{}, err
	}
	n, err = fmt.Sscanf(viewStr, "%f%s", &baseNumber, &suffix)
	if err != nil {
		return conf.Bangumi{}, err
	}
	if n != 2 {
		return conf.Bangumi{}, errors.New(viewStr + " type error")
	}
	if suffix == "亿次播放" {
		bangumi.View = int64(baseNumber * 100000000)
	} else if suffix == "万次播放" {
		bangumi.View = int64(baseNumber * 10000)
	} else if suffix == "次播放" {
		bangumi.View = int64(baseNumber)
	} else {
		return conf.Bangumi{}, errors.New("unknown suffix: " + suffix)
	}

	// get score
	scoreStr, err := util.JsonGetStr(dataDict, "score")
	if err == nil {
		var score float64
		n, err = fmt.Sscanf(scoreStr, "%f", &score)
		if err != nil {
			return conf.Bangumi{}, err
		}
		if n != 1 {
			return conf.Bangumi{}, errors.New(scoreStr + " type error")
		}
		bangumi.Score = score
	} else {
		bangumi.Score = -1
	}

	return bangumi, nil
}
