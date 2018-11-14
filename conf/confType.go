package conf

type Conf struct {
	MaxGoroutineNum int
	DB              DBConf `json:"database"`
	Network         NetworkConf
	StartAid        int
	EndAid          int
}

type DBConf struct {
	User            string
	Passwd          string
	Host            string
	DBname          string
	MaxOpenConn     int
	MaxIdleConn     int
	MaxConnLifeTime int
}

type NetworkConf struct {
	UseProxy   bool
	UserAgent  string
	RetryTimes int
}
