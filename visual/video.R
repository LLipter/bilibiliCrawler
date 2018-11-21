library(DBI)
library(RMySQL)
library(ggplot2)
library(scales)
Sys.setlocale(locale="UTF-8")

# connect to database
con <- dbConnect(MySQL(), host="localhost", dbname="bilibili", user="root", password="57575207")

video.rawdata <- dbGetQuery(con, "SELECT * FROM video WHERE status=0 AND aid%10000=0");

# pie chart
class <- c("<1k", "1k-10k", "10k-100k", "100k-1M", ">1M")
counts <- NULL
counts <- append(counts, lengths(subset(video.rawdata,view<1000))[1])
counts <- append(counts, lengths(subset(video.rawdata,view>=1000 & view<10000))[1])
counts <- append(counts, lengths(subset(video.rawdata,view>=10000 & view<100000))[1])
counts <- append(counts, lengths(subset(video.rawdata,view>=100000 & view<1000000))[1])
counts <- append(counts, lengths(subset(video.rawdata,view>=1000000))[1])
video.pie.data <- data.frame(class, counts)
video.pie.data.cumsum <- cumsum(counts)
video.pie.data.sum <- sum(counts)

video.pie <- ggplot(video.pie.data, aes(x="", y=counts, fill=class))
video.pie <- video.pie + geom_bar(width=1, stat = "identity")
video.pie <- video.pie + coord_polar("y", start=0)
video.pie <- video.pie + labs(title="View Distribution")
video.pie <- video.pie + theme(plot.title=element_text(hjust=0.5))
video.pie <- video.pie + theme(axis.ticks=element_blank())
video.pie <- video.pie + theme(legend.title=element_blank())
video.pie <- video.pie + theme(axis.title.x=element_blank())
video.pie <- video.pie + theme(axis.title.y=element_blank())
video.pie <- video.pie + theme(axis.text.x=element_blank())
video.pie <- video.pie + theme(axis.text.y=element_blank())
video.pie <- video.pie + theme(panel.border=element_blank())
video.pie <- video.pie + theme(panel.grid=element_blank())
video.pie <- video.pie + scale_fill_brewer("Blues", limits=class, direction=-1)
video.pie <- video.pie + geom_text(aes(y=counts[1]/2, label=percent(counts[1]/video.pie.data.sum)), size=5)
video.pie <- video.pie + geom_text(aes(y=counts[2]/2 + video.pie.data.cumsum[1], label=percent(counts[2]/video.pie.data.sum)), size=5)
video.pie <- video.pie + geom_text(aes(y=counts[3]/2 + video.pie.data.cumsum[2], label=percent(counts[3]/video.pie.data.sum)), size=5)
ggsave("assets/view_distribution.png")

a <- cumsum(video.pie.data["counts"])

b <- percent(unlist(video.pie.data["counts"]/100))

# close database
dbDisconnect(con)

  
