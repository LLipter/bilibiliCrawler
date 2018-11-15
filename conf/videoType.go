package conf

import "time"

type Video struct {
	Status int64

	Aid       int64 `json:"aid"`
	View      int64 `json:"view"`
	Dannmaku  int64 `json:"danmaku"`
	Reply     int64 `json:"reply"`
	Favorite  int64 `json:"favorite"`
	Coin      int64 `json:"coin"`
	Share     int64 `json:"share"`
	His_rank  int64 `json:"his_rank"`
	Support   int64 `json:"like"`
	Dislike   int64 `json:"dislike"`
	Copyright int64 `json:"copyright"`

	Cid     int64     `json:"cid"`
	Tid     int64     `json:"tid"`
	Title   string    `json:"title"`
	Pubdate time.Time `json:"created"`
	Mid     int64     `json:"mid"`
}

type Info struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    Video  `json:"data"`
}
