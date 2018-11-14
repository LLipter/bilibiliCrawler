# Overall

This is a simple crawler design for [bilibili](www.bilibili.com) that can collect basic data associate with each video. For example, it will collect the number of `comment`, `dannmaku`, `coin` of each video.

# Performance
After tuning parameters in configuration file, it can collect approximately 100k video data per hour. However, this number is trivial compared with the total video number in [bilibili](www.bilibili.com). By the time 2018/11/14, the max `av` number is larger than 35 million, which means it will cost this program  half a month to collect all video data, hopefully.

# How to run it

### Database

create databsase and tables with the following command

`mysql> source table.sql`

### Configuration file

rename `config-default.json` to `config.json` and edit it.

### Proxy

If you want to use proxy, change `GetProxy()` function in proxy/proxy.go to provide your proxy ip address.

### Compile codes

`go build`

### Start run it!

`./bilibiliVideoDataCrawler`
