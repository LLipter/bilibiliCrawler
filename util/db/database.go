package db

import (
	"database/sql"
	"fmt"
	"github.com/LLipter/bilibiliCrawler/conf"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var (
	connPool *sql.DB // database connection poold
)

func init() {
	var err error
	connPool, err = sql.Open("mysql", conf.DBConnStr)
	if err != nil {
		fmt.Printf("cannot create database connection pool, %v\n", err)
		os.Exit(1)
	}
	err = connPool.Ping()
	if err != nil {
		fmt.Printf("cannot access database, %v\n", err)
		os.Exit(1)
	}
	connPool.SetMaxOpenConns(conf.DBconfig.MaxOpenConn)
	connPool.SetConnMaxLifetime(time.Second * 9)

}

func CloseDatabase() {
	if connPool != nil {
		connPool.Close()
	}
}

func InsertVideo(video conf.Video) error {
	var pubdate interface{}
	if video.Pubdate.IsZero() {
		pubdate = nil
	} else {
		pubdate = video.Pubdate
	}

	_, err := connPool.Exec(
		"INSERT INTO video VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
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
		online.Data.Douga,
		online.Data.Anime,
		online.Data.Guochuang,
		online.Data.Music,
		online.Data.Dance,
		online.Data.Game,
		online.Data.Technology,
		online.Data.Life,
		online.Data.Kichiku,
		online.Data.Fashion,
		online.Data.Ad,
		online.Data.Ent,
		online.Data.Cinephile,
		online.Data.Cinema,
		online.Data.Tv,
		online.Data.Movie,
		online.AllCount,
		online.WebOnline,
		online.PlayOnline,
	)
	if err != nil {
		return err
	}

	return nil
}
