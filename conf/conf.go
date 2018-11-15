package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/LLipter/bilibiliVideoDataCrawler/daemon"
	"io/ioutil"
	"os"
)

var (
	DBConnStr       string
	MaxOpenConn     int
	RetryTimes      int
	UseProxy        bool
	UserAgent       string
	MaxCrawlerNum   int
	StartAid        int
	EndAid          int
)

func init() {
	buf, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Printf("cannot open configuration file, %v", err)
		os.Exit(1)
	}
	var config Conf
	err = json.Unmarshal(buf, &config)
	if err != nil {
		fmt.Printf("invalid configuration file, %v", err)
		os.Exit(1)
	}

	// get database connection string
	var strBuf bytes.Buffer
	strBuf.WriteString(config.DB.User)
	strBuf.WriteString(":")
	strBuf.WriteString(config.DB.Passwd)
	strBuf.WriteString("@tcp(")
	strBuf.WriteString(config.DB.Host)
	strBuf.WriteString(")/")
	strBuf.WriteString(config.DB.DBname)
	strBuf.WriteString("?charset=utf8")
	DBConnStr = strBuf.String()

	// get database connection parameters
	MaxOpenConn = config.DB.MaxOpenConn

	// get network configuration
	RetryTimes = config.Network.RetryTimes
	UseProxy = config.Network.UseProxy
	UserAgent = config.Network.UserAgent

	// get go routine max number
	MaxCrawlerNum = config.MaxCrawlerNum

	// get end aid
	EndAid = config.EndAid

	// check whether run as daemon
	if config.IsDaemon {
		daemon.Daemonize()
	}

}
