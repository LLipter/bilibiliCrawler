package conf

import "time"

type OnlineJson struct {
	Code       int64       `json:"code"`
	Message    string      `json:"message"`
	Data       RegionCount `json:"data"`
	AllCount   int64       `json:"all_count"`  // 最新投稿
	WebOnline  int64       `json:"web_online"` // 在线人数
	PlayOnline int64       `json:"play_online"`

	Timestamp time.Time
}

type RegionCount struct {
	Douga      int64 `json:"1"`   // 动画
	Anime      int64 `json:"13"`  // 番剧
	Guochuang  int64 `json:"167"` // 国创
	Music      int64 `json:"3"`   // 音乐
	Dance      int64 `json:"129"` // 舞蹈
	Game       int64 `json:"4"`   // 游戏
	Technology int64 `json:"36"`  // 科技
	Life       int64 `json:"160"` // 生活
	Kichiku    int64 `json:"119"` // 鬼畜
	Fashion    int64 `json:"155"` // 时尚
	Ad         int64 `json:"165"` // 广告
	Ent        int64 `json:"5"`   // 娱乐
	Cinephile  int64 `json:"181"` // 影视
	Cinema     int64 `json:"177"` // 放映厅

	Tv    int64 `json:"11"` // 电视剧
	Movie int64 `json:"23"` // 电影
}
