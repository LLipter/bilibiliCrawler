package conf

import "time"

type Video struct {
	Aid        int64
	Status     int64
	Title      string
	Pubdate    time.Time
	Owner      int64
	Duration   int64
	View       int64
	Dannmaku   int64 `json:"danmaku"`
	Reply      int64
	Favorite   int64
	Coin       int64
	Share      int64
	Now_rank   int64
	His_rank   int64
	Support    int64
	Dislike    int64
	No_reprint int64
	Copyright  int64
	Pages      []Page
}

type Info struct {
	Code    int64
	Message string
	Data    Video
}

type Page struct {
	PageNo   int64
	Chatid   int64
	Duration int64
	Subtitle string
}
