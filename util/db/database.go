package db

import (
	"database/sql"
	"fmt"
	"github.com/LLipter/bilibiliCrawler/conf"
	_ "github.com/go-sql-driver/mysql"
	"os"
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
