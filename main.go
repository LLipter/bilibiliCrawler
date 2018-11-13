package main

import (
	"fmt"
	_ "github.com/LLipter/bilibili-report/conf"
	"github.com/LLipter/bilibili-report/crawler"
	"github.com/LLipter/bilibili-report/util/db"
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

	fmt.Println("init successfully. ")
}

func cleanup() {
	if logFile != nil {
		logFile.Close()
	}
	db.CloseDatabase()
}

func main() {
	crawler.CrawlVideo(1001, 2000)
	cleanup()
}
