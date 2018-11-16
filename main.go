package main

import (
	"fmt"
	"github.com/LLipter/bilibiliCrawler/conf"
	"github.com/LLipter/bilibiliCrawler/crawler"
	"github.com/LLipter/bilibiliCrawler/proxy"
	"github.com/LLipter/bilibiliCrawler/util/db"
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
		log.Fatalln(fmt.Sprintf("cannot open log file, %v", err))
	}
	log.SetOutput(logFile)

	// start proxy
	if conf.NetworkConfig.UseProxy {
		err = proxy.GetProxy()
		if err != nil {
			log.Fatalln(err)
		}
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
	log.Println("begin crawling")
	if os.Args[1] == "-v" {
		crawler.CrawlVideo(conf.VideoCrawlerConfig.StartAid, conf.VideoCrawlerConfig.EndAid)
	} else if os.Args[1] == "-o" {
		crawler.CrawOnline()
	}
	log.Println("end crawling")
}
