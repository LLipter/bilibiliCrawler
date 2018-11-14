# Overall

This is a simple crawler design for [bilibili](www.bilibili.com) that can collect basic data associate with each video. For example, it will collect the number of `comment`, `dannmaku`, `coin` of each video.

# Performance
After tuning parameters in configuration file and run it in my server, it can collect approximately 100k video data per hour. However, this number is trivial compared with the total video number in [bilibili](www.bilibili.com). By the time 2018/11/14, the max `av` number is larger than 35 million, which means it will cost this program half a month to collect all video data, hopefully.

### Bottleneck

The speed of Internet.

# How to run it

### Database

create databsase and tables with the following command

`mysql> source table.sql`

### Configuration file

rename `config-default.json` to `config.json` and edit it.

### Proxy

If you want to use proxy, change `GetProxy()` function in `proxy/proxy.go` to provide your proxy ip addresses. **The quality of your proxy server has a great influence on the overall performance**

### Compile codes

`go build`

### Start running it!

`./bilibiliVideoDataCrawler`

# Tuning parameters in configuration file

### maxCrawlerNum

This parameter determine the max number of crawler go routines. Although a go routine is much lighter than a thread, this number cannot be **too large**. Because in many operator system, there's a limit on how many file descriptor one process can use, **especially in MacOS**. Since a socket will consume a file descriptor, if this number is too large, the following error may occur. In my server, `500` seems reasonable. You may use `ulimit -a` to check the limit of file descriptor in your system.

`proxyconnect tcp: dial tcp 112.74.41.42:80: socket: too many open files`

### retryTimes

This parameter determine the max number a go routine failed to crawl a video data until it gives up. It highly depends on your internet reliability and your proxy server's reliability if you choose to use proxy. No one knows why some data isn't properly transmitted. If you're not confident with your internet, try a higher value. In my case, I use `10`, `20` or even `50`.

### maxOpenConn

This parameter determine the max number of connection that mysql connection pool holds. Mysql server limits the number of connection to `100` by default, so I set this parameter to `90` in case some emergency may occur.
