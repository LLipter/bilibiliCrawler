package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/LLipter/bilibiliCrawler/daemon"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	DBconfig       DBConf
	NetworkConfig  NetworkConf
	DBConnStr      string
	isDaemon       = false
	IsCrawlVideo   = false
	IsCrawlOnline  = false
	IsCrawlBangumi = false
	StartAid       int
	EndAid         int
)

func init() {
	if !isValidParameter() {
		usage()
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

	// get database connection string
	var strBuf bytes.Buffer
	strBuf.WriteString(config.DB.User)
	strBuf.WriteString(":")
	strBuf.WriteString(config.DB.Passwd)
	strBuf.WriteString("@tcp(")
	strBuf.WriteString(config.DB.Host)
	strBuf.WriteString(")/")
	strBuf.WriteString(config.DB.DBname)
	strBuf.WriteString("?charset=utf8&parseTime=true&loc=Local")
	DBConnStr = strBuf.String()

	// check whether run as daemon
	if isDaemon {
		daemon.Daemonize()
	}

}

func usage() {
	fmt.Println("usage: bilibiliCrawler [-v[d] startAid endAid][-o[d]][-b[d]]")
	fmt.Println("   -v: crawl video data")
	fmt.Println("   -o: crawl online data")
	fmt.Println("   -b: crawl bangumi data")
	fmt.Println("   -d: run as daemon process")

	os.Exit(1)
}

func isValidParameter() bool {
	if len(os.Args) < 2 {
		return false
	}
	var err error
	arg := os.Args[1]
	if arg == "-v" {
		if len(os.Args) < 4 {
			return false
		}
		StartAid, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return false
		}
		EndAid, err = strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println(err)
			return false
		}
		IsCrawlVideo = true
		return true
	} else if arg == "-o" {
		IsCrawlOnline = true
		return true
	} else if arg == "-b" {
		IsCrawlBangumi = true
		return true
	} else if arg == "-vd" {
		if len(os.Args) < 4 {
			return false
		}
		StartAid, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return false
		}
		EndAid, err = strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println(err)
			return false
		}
		IsCrawlVideo = true
		isDaemon = true
		return true
	} else if arg == "-od" {
		IsCrawlOnline = true
		isDaemon = true
		return true
	} else if arg == "-bd" {
		IsCrawlBangumi = true
		isDaemon = true
		return true
	} else {
		fmt.Println("unknown parameter")
	}

	return false
}
