package conf

import "time"

var (
	RetryTimes       = 3
	MaxGoroutinueNum = 200
	UseProxy         = false
	UserAgent        = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"
	MaxOpenConn      = 100
	MaxIdleConn      = 30
	MaxConnLifeTime  = time.Minute * 10
)
