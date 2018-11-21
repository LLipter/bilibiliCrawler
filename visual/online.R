library(DBI)
library(RMySQL)
library(ggplot2)
library(scales)
Sys.setlocale(locale="UTF-8")


# connect to database
con <- dbConnect(MySQL(), host="localhost", dbname="bilibili", user="root", password="57575207")

online.rawdata <- dbGetQuery(con, "SELECT * FROM online");
online.data <- data.frame(ts=as.POSIXct(online.rawdata$ts),
                        allcount=online.rawdata$allcount,
                        webonline=online.rawdata$webonline)

oo <- options(scipen=200)
online.plot <- ggplot(data=online.data)
online.plot <- online.plot + geom_line(mapping=aes(x=ts, y=webonline, colour="Online User"))
online.plot <- online.plot + geom_line(mapping=aes(x=ts, y=allcount*30, colour="New Video"))
# online.plot <- online.plot + scale_y_continuous(sec.axis=sec_axis(~.+1000000, name="allcount"))
online.plot <- online.plot + labs(title="Bilibili Daily Active User")
online.plot <- online.plot + ylab("Number of online users")
online.plot <- online.plot + xlab("")
online.plot <- online.plot + scale_y_continuous(sec.axis=sec_axis(~./30, name="Number of New Video"))
online.plot <- online.plot + scale_x_datetime(breaks=date_breaks("6 hour"), date_labels="%b %d %H:00")
online.plot <- online.plot + theme(axis.text.x=element_text(angle=45,vjust=1,hjust=1))


online.plot <- online.plot + theme(plot.title=element_text(hjust = 0.5))
ggsave("assets/daily_active_user.png", width=14)
options(oo)

# close database
dbDisconnect(con)
