library(DBI)
library(RMySQL)
library(ggplot2)
Sys.setlocale(locale="UTF-8")

# connect to database
con <- dbConnect(MySQL(), host="localhost", dbname="bilibili", user="root", password="57575207")

video.rawdata <- dbGetQuery(con, "SELECT * FROM video WHERE status=0 AND aid%10000=0");

# close database
dbDisconnect(con)
