# Overall

This is a simple crawler design for [bilibili](www.bilibili.com) that can collect basic data associate with each video. For example, it will collect the number of `comment`, `dannmaku`, `coin`, etc.

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
