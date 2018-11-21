library(DBI)
library(RMySQL)
library(ggplot2)
library(scales)
Sys.setlocale(locale="UTF-8")

# connect to database
con <- dbConnect(MySQL(), host="localhost", dbname="bilibili", user="root", password="57575207")
video.rawdata <- dbGetQuery(con, "SELECT * FROM video WHERE status=0 LIMIT 500000");

# view chart
class <- c("<1k", "1k-10k", "10k-100k", "100k-1M", ">1M")
counts <- NULL
counts <- append(counts, lengths(subset(video.rawdata,view<1000))[1])
counts <- append(counts, lengths(subset(video.rawdata,view>=1000 & view<10000))[1])
counts <- append(counts, lengths(subset(video.rawdata,view>=10000 & view<100000))[1])
counts <- append(counts, lengths(subset(video.rawdata,view>=100000 & view<1000000))[1])
counts <- append(counts, lengths(subset(video.rawdata,view>=1000000))[1])
video.view.data <- data.frame(class, counts)
video.view.data.cumsum <- cumsum(counts)
video.view.data.sum <- sum(counts)

video.view <- ggplot(video.view.data, aes(x="", y=counts, fill=class))
video.view <- video.view + geom_bar(width=1, stat = "identity")
video.view <- video.view + coord_polar("y", start=0)
video.view <- video.view + labs(title="View Distribution")
video.view <- video.view + theme(plot.title=element_text(hjust=0.5))
video.view <- video.view + theme(axis.ticks=element_blank())
video.view <- video.view + theme(legend.title=element_blank())
video.view <- video.view + theme(axis.title.x=element_blank())
video.view <- video.view + theme(axis.title.y=element_blank())
video.view <- video.view + theme(axis.text.x=element_blank())
video.view <- video.view + theme(axis.text.y=element_blank())
video.view <- video.view + theme(panel.border=element_blank())
video.view <- video.view + theme(panel.grid=element_blank())
video.view <- video.view + scale_fill_brewer("Blues", limits=class, direction=-1)
video.view <- video.view + geom_text(aes(y=counts[1]/2, label=percent(counts[1]/video.view.data.sum)), size=5)
video.view <- video.view + geom_text(aes(y=counts[2]/2 + video.view.data.cumsum[1], label=percent(counts[2]/video.view.data.sum)), size=5)
video.view <- video.view + geom_text(aes(y=counts[3]/2 + video.view.data.cumsum[2], label=percent(counts[3]/video.view.data.sum)), size=5)
video.view <- video.view + geom_text(aes(y=counts[4]/2 + video.view.data.cumsum[3], label=percent(counts[4]/video.view.data.sum), x=1.1), size=5)
video.view <- video.view + geom_text(aes(y=counts[5]/2 + video.view.data.cumsum[4], label=percent(counts[5]/video.view.data.sum), x=1.3), size=5)
ggsave("assets/view_distribution.png")

# pubdate chart
video.pubdate.data <- as.POSIXct(video.rawdata$pubdate)
video.pubdate.data.bymonth <- cut(video.pubdate.data, breaks="month")
video.pubdate.data.split <- split(video.pubdate.data, video.pubdate.data.bymonth)
video.pubdate.cnt <- data.frame(pubdate=as.POSIXct(names(video.pubdate.data.split)), counts=lengths(video.pubdate.data.split))

oo <- options(scipen=200)
video.pubdate <- ggplot(data=video.pubdate.cnt)
video.pubdate <- video.pubdate + geom_line(mapping=aes(x=pubdate, y=counts))
video.pubdate <- video.pubdate + labs(title="Bilibili Monthly New Video")
video.pubdate <- video.pubdate + ylab("Number of New Video")
video.pubdate <- video.pubdate + xlab("")
video.pubdate <- video.pubdate + scale_x_datetime(breaks=date_breaks("3 month"), date_labels="%Y %b")
video.pubdate <- video.pubdate + theme(axis.text.x=element_text(angle=45,vjust=1,hjust=1))
video.pubdate <- video.pubdate + theme(plot.title=element_text(hjust=0.5))
ggsave("assets/monthly_new_video.png", width=14, height=7)
options(oo)

# close database
dbDisconnect(con)

  
