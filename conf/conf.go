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
	DBconfig           DBConf
	NetworkConfig      NetworkConf
	VideoCrawlerConfig VideoCrawlerConf

	DBConnStr string
)

func init() {
	if len(os.Args) != 2 {
		usage()
	}

	if isValidParameter(os.Args[1]) {
		fmt.Println("unknown parameter")
		usage()
		os.Exit(1)
	}

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

	// get configuration
	DBconfig = config.DB
	NetworkConfig = config.Network
	VideoCrawlerConfig = config.VideoCrawler

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

	// check whether run as daemon
	if os.Args[1] == "-v" && config.VideoCrawler.IsDaemon {
		daemon.Daemonize()
	}

}

func usage() {
	fmt.Println("usage: bilibiliCrawler [-v]")
	os.Exit(1)
}

func isValidParameter(arg string) bool {
	if arg == "-v" {
		return true
	}
	return false
}
