package util

import "time"

type Video struct {
	Aid        int
	Status     int
	Title	   string
	Pubdate	   time.Time
	Duration   int
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

type DBConf struct {
	User   string
	Passwd string
	Host   string
	Dbname string
}

type Page struct{
	Chatid	int
	Duration int
	Subtitle string
	PageNo int
}
