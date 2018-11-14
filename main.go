package main

import (
	"fmt"
	"github.com/LLipter/bilibiliVideoDataCrawler/conf"
	"github.com/LLipter/bilibiliVideoDataCrawler/crawler"
	"github.com/LLipter/bilibiliVideoDataCrawler/proxy"
	"github.com/LLipter/bilibiliVideoDataCrawler/util/db"
	"log"
	"os"
)

var (
	logFile *os.File
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var err error
	logFile, err = os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot open log file, %v", err))
	}
	log.SetOutput(logFile)

	// start proxy
	if conf.UseProxy {
		proxy.GetProxy()
		go proxy.GetProxyRoutine()
	}

	fmt.Println("init successfully. ")
}

func cleanup() {
	if logFile != nil {
		logFile.Close()
	}
	db.CloseDatabase()
}

func main() {
	defer cleanup()
	crawler.CrawlVideo(conf.StartAid, conf.EndAid)
}
