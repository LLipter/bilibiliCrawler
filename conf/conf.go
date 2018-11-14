package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var (
	DBConnStr       string
	MaxOpenConn     int
	MaxIdleConn     int
	MaxConnLifeTime time.Duration
	RetryTimes      int
	UseProxy        bool
	UserAgent       string
	MaxGoroutineNum int
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
	MaxIdleConn = config.DB.MaxIdleConn
	MaxConnLifeTime = time.Duration(config.DB.MaxConnLifeTime) * time.Minute

	// get network configuration
	RetryTimes = config.Network.RetryTimes
	UseProxy = config.Network.UseProxy
	UserAgent = config.Network.UserAgent

	// get go routine max number
	MaxGoroutineNum = config.MaxGoroutineNum

	// get start aid and end aid
	StartAid = config.StartAid
	EndAid = config.EndAid
}
