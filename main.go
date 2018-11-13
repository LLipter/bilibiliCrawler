package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/LLipter/bilibili-report/util"
)

var (
	retryTimes			= 3
	maxGoroutinueNum	= 200
	useProxy			= false
	userAgent			= "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"
	wg               	sync.WaitGroup
	logFile          	*os.File
)

func getResp(addr string) (*http.Response,error){
	client := http.Client{}
	req, err := http.NewRequest("GET", addr, nil)
	req.Header.Add("User-Agent", userAgent)
	if useProxy {
		// TODO: add proxies pool
		urlproxy, err := url.Parse("http://183.245.99.52:80")
		if err != nil {
			return nil, err
		}
		client = http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlproxy),
			},
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp,nil
}

func getVideoBasicData(aid int) (util.Info, error){
	addr := "https://api.bilibili.com/archive_stat/stat?aid="+strconv.Itoa(aid)
	resp, err := getResp(addr)
	if err != nil {
		return util.Info{}, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return util.Info{}, err
	}

	var info util.Info
	err = json.Unmarshal(data, &info)
	if err != nil {
		if strings.HasPrefix(err.Error(), "json: cannot unmarshal string into Go struct field"){
			return util.Info{Code:-1},nil
		}
		return util.Info{}, err
	}

	return info, nil
}

func getVideoPostTime(aid int, video *util.Video) error{
	addr := "https://www.bilibili.com/video/av"+strconv.Itoa(aid)
	resp, err := getResp(addr)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	var jsonstr string
	doc.Find("script").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		content := selection.Text()
		if(strings.HasPrefix(content, "window.__INITIAL_STATE__=")){
			jsonstr = content[25:len(content)-122]
			return false
		}
		return true
	})

	var jsonObj map[string]interface{}
	json.Unmarshal([]byte(jsonstr), &jsonObj)

	util.PrintJson(jsonObj)

	videoData, ok := jsonObj["videoData"]
	if !ok {
		return errors.New("missing 'videoData'")
	}
	videoDataMap, ok := videoData.(map[string]interface{})
	if !ok {
		return errors.New("'videoData' type error")
	}

	// retrieve title
	title, ok := videoDataMap["title"]
	if !ok {
		return  errors.New("missing 'title'")
	}
	titleStr, ok := title.(string)
	if !ok {
		return  errors.New("'title' type error")
	}
	video.Title = titleStr

	// retrieve pubdate
	pubdate, ok := videoDataMap["pubdate"]
	if !ok {
		return errors.New("missing 'pubdate'")
	}
	pubdateFloat, ok := pubdate.(float64)
	if !ok {
		return errors.New("'pubdate' type error")
	}
	video.Pubdate = time.Unix(int64(pubdateFloat), 0)

	// retrieve duration
	duration, ok := videoDataMap["duration"]
	if !ok {
		return errors.New("missing 'duration'")
	}
	durationFloat, ok := duration.(float64)
	if !ok {
		return errors.New("'duration' type error")
	}
	video.Duration = int(durationFloat)

	// retrieve pages
	pages, ok := videoDataMap["pages"]
	if !ok {
		return errors.New("missing 'pages'")
	}
	pagesArray, ok := pages.([]interface {})
	if !ok {
		return errors.New("'pages' type error")
	}
	for page := range pagesArray{
		
	}
	fmt.Printf("%T\n", pages)





	return nil
}


func getVideoData(aid int) (util.Info, error) {
	data, err := getVideoBasicData(aid)
	if err != nil {
		return data, err
	}


	return data, nil
}

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

func crawler(aid int) error {
	info, err := getVideoData(aid)
	if err != nil {
		return err
	}

	var video util.Video

	if info.Code != 0 {
		video.Status = 1
		video.Aid = aid
	} else {
		video = info.Data
	}

	err = util.InsertVideo(video)
	if err != nil {
		return err
	}

	return nil
}

func crawlerRoutine(aid int) {
	defer wg.Done()
	for t := 0; t < retryTimes; t++ {
		err := crawler(aid)
		if err == nil {
			return
		} else {
			log.Printf("aid=%d crawler failed, %v\n", aid, err)
		}
	}
	// failed with unknown reason
	var video util.Video
	video.Status = 2
	video.Aid = aid
	err := util.InsertVideo(video)
	if err != nil {
		log.Printf("aid=%d insertion failed, %v\n", aid, err)
	}
}

func cleanup() {
	if logFile != nil {
		logFile.Close()
	}
	util.CloseDatabase()
}

func main() {

	//for i := 1; i <= 300; i++ {
	//	for runtime.NumGoroutine() > maxGoroutinueNum {
	//		time.Sleep(time.Second)
	//	}
	//	wg.Add(1)
	//	go crawlerRoutine(i)
	//}

	var v util.Video
	getVideoPostTime(35679613, &v)
	fmt.Println(v.Title)
	fmt.Println(v.Pubdate)
	fmt.Println(v.Duration)


	wg.Wait()

	cleanup()

}
