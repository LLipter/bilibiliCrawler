# Overall

This is a simple crawler design for [bilibili](www.bilibili.com) that can collect basic data associate with each video. For example, it will collect the number of `comment`, `dannmaku`, `coin` of each video.

# Performance
After tuning parameters in configuration file and run it in my server, it can collect approximately 300k video data per hour. However, this number is trivial compared with the total video number in [bilibili](www.bilibili.com). By the time 2018/11/14, the max `av` number is larger than 35 million, which means it will cost this program several days to collect all video data, hopefully.

### Bottleneck

The speed of Internet. 

As you can see, outward internet flow reaches upper limit as soon as I running this program

![](assets/monitor.png)

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

`./bilibiliCrawler [-vo]`


# Tuning parameters in configuration file

### maxCrawlerNum

This parameter determine the max number of crawler go routines. Although a go routine is much lighter than a thread, this number cannot be **too large**. Because the bottleneck is the speed of internet. In my server, `300` seems reasonable. 

### retryTimes

This parameter determine the max number a go routine failed to crawl a video data until it gives up. It highly depends on your internet reliability and your proxy server's reliability if you choose to use proxy. No one knows why some data isn't properly transmitted. If you're not confident with your internet, try a higher value. In my case, I use `10`, `20` or even `30`.

### maxOpenConn

This parameter determine the max number of connection that mysql connection pool holds. Mysql server limits the number of connection to `100` by default, so I set this parameter to `95` in case some emergency may occur.
