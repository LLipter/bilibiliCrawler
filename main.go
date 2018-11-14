package main

import (
	"fmt"
	_ "github.com/LLipter/bilibili-report/daemon"
	"github.com/LLipter/bilibili-report/conf"
	"github.com/LLipter/bilibili-report/crawler"
	"github.com/LLipter/bilibili-report/proxy"
	"github.com/LLipter/bilibili-report/util/db"
	"log"
	"os"
	"time"
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
		go proxy.GetProxies()
	}
	// make sure proxies is available
	time.Sleep(time.Second)

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
