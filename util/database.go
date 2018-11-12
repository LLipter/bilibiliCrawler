package util

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var(
	connPool		*sql.DB // database connection pool
	maxOpenConn		= 100
	maxIdleConn 	= 30
	maxConnLifeTime = time.Minute * 10
)

func init(){
	// change to your database configuration file path, see dbconfig-sample.json
	connStr, err := LoadDBConf("dbconfig.json")
	if err != nil {
		fmt.Printf("cannot open database configuration file, %v", err)
		os.Exit(1)
	}
	connPool, err = sql.Open("mysql", connStr)
	if err != nil {
		fmt.Printf("cannot create database connection pool, %v", err)
		os.Exit(1)
	}
	err = connPool.Ping()
	if err != nil {
		fmt.Printf("cannot access database, %v", err)
		os.Exit(1)
	}

	connPool.SetMaxOpenConns(maxOpenConn)
	connPool.SetMaxIdleConns(maxIdleConn)
	connPool.SetConnMaxLifetime(maxConnLifeTime)
}

func CloseDatabase(){
	if connPool != nil{
		connPool.Close()
	}
}


func InsertVideo(video Video) error{
	stmt, err := connPool.Prepare("INSERT INTO video VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		video.Aid,
		video.Status,
		video.View,
		video.Dannmaku,
		video.Reply,
		video.Favorite,
		video.Coin,
		video.Share,
		video.Now_rank,
		video.His_rank,
		video.Support,
		video.Dislike,
		video.No_reprint,
		video.Copyright)
	if err != nil {
		return err
	}
	return nil
}