package util

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var(
	db	*sql.DB		// database connection pool
)

func init(){
	// change to your database configuration file path, see dbconfig-sample.json
	connStr, err := LoadDBConf("dbconfig.json")
	fmt.Println(connStr)
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDatabase(){
	if db != nil{
		db.Close()
	}
}


func InsertVideo(video Video) error{
	stmt, err := db.Prepare("INSERT INTO video VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}
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