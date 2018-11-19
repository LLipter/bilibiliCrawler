library(DBI)
library(RMySQL)
Sys.setlocale(locale="UTF-8")
con <- dbConnect(MySQL(), host="www.irran.top", dbname="bilibili", user="root", password="5720")

online.data <- dbGetQuery(con, "SELECT * FROM online order by ts LIMIT 100");




dbDisconnect(con)