package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LLipter/bilibiliCrawler/conf"
	"github.com/LLipter/bilibiliCrawler/util"
	"strconv"
	"time"
)

func CrawlBangumi() error {
	addr := "https://bangumi.bilibili.com/media/web_api/search/result?season_version=1&area=2&is_finish=1&copyright=-1&season_status=-1&season_month=-1&pub_date=-1&style_id=-1&order=3&st=1&sort=0&page=1&season_type=1&pagesize=20"
	buf, err := getResp(addr)
	if err != nil {
		return errors.New("cannot get total number: " + err.Error())
	}

	var bangumiJson conf.BangumiJson
	err = json.Unmarshal(buf, &bangumiJson)
	if err != nil {
		return errors.New("cannot get total number: " + err.Error())
	}

	pageSize := 20
	totalNo := int(bangumiJson.Data.Page.Total)
	pageNo := totalNo / pageSize
	if totalNo%pageSize != 0 {
		pageNo++
	}

	for page := 1; page <= 1; page++ {
		addr := "https://bangumi.bilibili.com/media/web_api/search/result?season_version=1&area=2&is_finish=1&copyright=-1&season_status=-1&season_month=-1&pub_date=-1&style_id=-1&order=3&st=1&sort=0&season_type=1&pagesize=20&page=" + strconv.Itoa(page)
		buf, err := getResp(addr)
		if err != nil {
			return errors.New(fmt.Sprintf("cannot get page %d: ", page) + err.Error())
		}

		var bangumiJson conf.BangumiJson
		err = json.Unmarshal(buf, &bangumiJson)
		if err != nil {
			return errors.New(fmt.Sprintf("cannot get page %d: ", page) + err.Error())
		}

		for _, data := range bangumiJson.Data.Bangumis {

			dataDict, ok := data.(map[string]interface{})
			if !ok {
				return pageError(page, err)
			}

			var bangumi conf.Bangumi

			// get media_id
			bangumi.MediaID, err = util.JsonGetInt64(dataDict, "media_id")
			if err != nil {
				return pageError(page, err)
			}

			// get season id
			bangumi.Sid, err = util.JsonGetInt64(dataDict, "season_id")
			if err != nil {
				return pageError(page, err)
			}

			// get season title
			bangumi.Title, err = util.JsonGetStr(dataDict, "title")
			if err != nil {
				return pageError(page, err)
			}

			dataDict, err = util.JsonGetDict(dataDict, "order")
			if err != nil {
				return pageError(page, err)
			}

			// get pubdate
			pubdate, err := util.JsonGetInt64(dataDict, "pub_date")
			if err != nil {
				return pageError(page, err)
			}
			bangumi.Pubdate = time.Unix(pubdate, 0)

			// get follow
			followStr, err := util.JsonGetStr(dataDict, "follow")
			if err != nil {
				return pageError(page, err)
			}
			var baseNumber float64
			var suffix string
			n, err := fmt.Sscanf(followStr, "%f%s", &baseNumber, &suffix)
			if err != nil {
				return pageError(page, err)
			}
			if n != 2 {
				return pageError(page, errors.New(followStr+" type error"))
			}
			if suffix == "万人追番" {
				bangumi.Follow = int64(baseNumber * 10000)
			} else if suffix == "人追番" {
				bangumi.Follow = int64(baseNumber)
			} else {
				return pageError(page, errors.New("unknown suffix: "+suffix))
			}

			// get view
			viewStr, err := util.JsonGetStr(dataDict, "play")
			if err != nil {
				return pageError(page, err)
			}
			n, err = fmt.Sscanf(viewStr, "%f%s", &baseNumber, &suffix)
			if err != nil {
				return pageError(page, err)
			}
			if n != 2 {
				return pageError(page, errors.New(viewStr+" type error"))
			}
			if suffix == "亿次播放" {
				bangumi.View = int64(baseNumber * 100000000)
			} else if suffix == "万次播放" {
				bangumi.View = int64(baseNumber * 10000)
			} else if suffix == "次播放" {
				bangumi.View = int64(baseNumber)
			} else {
				return pageError(page, errors.New("unknown suffix: "+suffix))
			}

			// get score
			scoreStr, err := util.JsonGetStr(dataDict, "score")
			if err == nil {
				var score float64
				n, err = fmt.Sscanf(scoreStr, "%f", &score)
				if err != nil {
					return pageError(page, err)
				}
				if n != 1 {
					return pageError(page, errors.New(scoreStr+" type error"))
				}
				bangumi.Score = score
			} else {
				bangumi.Score = -1
			}

			// get eplist
			err = getEplist(int(bangumi.Sid), &bangumi)
			if err != nil {
				return pageError(page, err)
			}

			util.PrintJson(bangumi)
		}

	}

	return nil
}

func getEplist(sid int, bangumi *conf.Bangumi) error {
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
		err = getVideoBasicData(aid, &videoJson)
		if err != nil {
			return errors.New(fmt.Sprintf("aid=%d cannot get view: ", aid) + err.Error())
		}
		eplistJson.Data[i].View = videoJson.Data.View
	}

	bangumi.Epno = int64(len(eplistJson.Data))
	bangumi.Eplist = eplistJson.Data
	return nil
}

func pageError(page int, err error) error {
	return errors.New(fmt.Sprintf("page %d: ", page) + err.Error())
}
