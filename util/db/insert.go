package db

import (
	"errors"
	"github.com/LLipter/bilibiliCrawler/conf"
)

func InsertVideo(video conf.Video) error {
	var pubdate interface{}
	if video.Pubdate.IsZero() {
		pubdate = nil
	} else {
		pubdate = video.Pubdate
	}

	_, err := connPool.Exec(
		"INSERT INTO video VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		video.Aid,
		video.Status,
		video.Title,
		pubdate,
		video.Mid,
		video.Cid,
		video.Tid,
		video.View,
		video.Dannmaku,
		video.Reply,
		video.Favorite,
		video.Coin,
		video.Share,
		video.His_rank,
		video.Support,
		video.Dislike,
		video.Copyright,
		video.Pages,
	)
	if err != nil {
		return err
	}

	return nil
}

func InsertOnline(online conf.OnlineJson) error {
	_, err := connPool.Exec(
		"INSERT INTO online VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		online.Timestamp,
		online.Data.RegionCount.Douga,
		online.Data.RegionCount.Anime,
		online.Data.RegionCount.Guochuang,
		online.Data.RegionCount.Music,
		online.Data.RegionCount.Dance,
		online.Data.RegionCount.Game,
		online.Data.RegionCount.Technology,
		online.Data.RegionCount.Life,
		online.Data.RegionCount.Kichiku,
		online.Data.RegionCount.Fashion,
		online.Data.RegionCount.Ad,
		online.Data.RegionCount.Ent,
		online.Data.RegionCount.Cinephile,
		online.Data.RegionCount.Cinema,
		online.Data.RegionCount.Tv,
		online.Data.RegionCount.Movie,
		online.Data.AllCount,
		online.Data.WebOnline,
		online.Data.PlayOnline,
	)
	if err != nil {
		return err
	}

	return nil
}

func InsertBangumi(bangumi conf.Bangumi) error {
	tx, err := connPool.Begin()
	if err != nil {
		return errors.New("transaction begin failed : " + err.Error())
	}

	// some bangumi doesn't have score
	var score interface{}
	if bangumi.Score < 0 {
		score = nil
	} else {
		score = bangumi.Score
	}

	_, err = connPool.Exec(
		"INSERT INTO bangumi VALUES(?, ?, ?, ?, ?, ?, ?, ?);",
		bangumi.Sid,
		bangumi.Title,
		bangumi.Pubdate,
		bangumi.Epno,
		score,
		bangumi.Follow,
		bangumi.View,
		bangumi.MediaID,
	)
	if err != nil {
		return rollback(tx, err)
	}

	for _, ep := range bangumi.Eplist {
		_, err = connPool.Exec(
			"INSERT INTO episode VALUES(?, ?, ?, ?, ?, ?);",
			bangumi.Sid,
			ep.Index,
			ep.Aid,
			ep.View,
			ep.Cid,
			ep.Epid,
		)
		if err != nil {
			return rollback(tx, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("transaction Commit failed : " + err.Error())
	}
	return nil
}
