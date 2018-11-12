package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/LLipter/bilibili-report/util"
	_ "github.com/go-sql-driver/mysql"
)

type Video struct {
	Aid        int
	Status     int
	View       int
	Dannmaku   int
	Reply      int
	Favorite   int
	Coin       int
	Share      int
	Now_rank   int
	His_rank   int
	Support    int
	Dislike    int
	No_reprint int
	Copyright  int
}

type Info struct {
	Code    int
	Message string
	Data    Video
}

func sendRequest(addr string) (Info, error) {
	//urlproxy, err := url.Parse("http://183.245.99.52:80")
	//if err != nil {
	//	return err
	//}
	//
	//client := http.Client{
	//	Transport: &http.Transport{
	//		Proxy: http.ProxyURL(urlproxy),
	//	},
	//}
	//resp, err := client.Get(addr)
	//if err != nil {
	//	return err
	//}

	resp, err := http.Get(addr)
	if err != nil {
		return Info{}, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Info{}, err
	}

	var info Info
	err = json.Unmarshal(data, &info)
	if err != nil {
		return Info{}, err
	}

	return info,nil
}

func Insert(video Video, db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO video VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(
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
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
}



func main() {
	info, err := sendRequest("https://api.bilibili.com/archive_stat/stat?aid=2")
	if err != nil {
		fmt.Println(err)
	}

	util.PrintJson(info)
	fmt.Println(info)


	db, err := sql.Open("mysql", "root:5720@tcp(www.irran.top:3306)/bilibili")
	defer db.Close()
	Insert(info.Data, db)


}
