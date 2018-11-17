package conf

import "time"

type BangumiJson struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    BangumiData `json:"result"`
}

type BangumiData struct {
	Bangumis []interface{} `json:"data"`
	Page     Page          `json:"page"`
}

type Page struct {
	Number int64 `json:"num"`
	Size   int64 `json:"size"`
	Total  int64 `json:"total"`
}

type Bangumi struct {
	Sid     int64
	Title   string
	Pubdate time.Time
	Epno    int64
	MediaID int64
	Score   float64
	Follow  int64
	View    int64
	Eplist  []Ep
}

type EplistJson struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    []Ep   `json:"result"`
}

type Ep struct {
	Aid   int64  `json:"avid"`
	Cid   int64  `json:"cid"`
	Epid  int64  `json:"episode_id"`
	Index string `json:"index"`
	View  int64
}
