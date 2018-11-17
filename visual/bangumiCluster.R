library(DBI)
library(RMySQL)
Sys.setlocale(locale="UTF-8")
con <- dbConnect(MySQL(), host="localhost", dbname="bilibili", user="root", password="57575207")
rawdata <- dbGetQuery(con, "SELECT * FROM bangumi WHERE ABS(view-view_calculated) < view*0.1 AND epno>10 ORDER BY view DESC LIMIT 50")
bangumidata <- rawdata[c("follow", "view_calculated")]
bangumidata <- scale(bangumidata)
row.names(bangumidata) <- rawdata$title

viewdata <- data.frame(view1=numeric(),
                        view3=numeric(),
                        view1n=numeric(),
                        view3n=numeric())    
for(sid in rawdata$sid){
    sqlstr <- paste("SELECT view FROM episode WHERE sid =", sid, "ORDER BY idx")
    epdata <- dbGetQuery(con, sqlstr)
    epdata <- epdata$view
    epdata <- scale(epdata)
    len <- length(epdata)
    row <- c(epdata[1],epdata[3],epdata[len],epdata[len-3])
    viewdata <- rbind(viewdata, row)
}
names(viewdata) <- c("view1","view3","view1n","view3n") 
bangumidata <- cbind(bangumidata, viewdata)

d <- dist(bangumidata)
fit.average <- hclust(d, method = "average")
plot(fit.average, hang = -1, cex = .5, main = "Average Linkage Clustering",
family='STXihei',
xlab="50 Most popular Anime in Bilibili",
ylab="",
sub="")

dbDisconnect(con)


    # library(NbClust)
    # set.seed(1234)
    # devAskNewPage(ask = TRUE)
    # nc <- NbClust(bangumidata, min.nc = 2, max.nc = 15, method = "kmeans")
    # table(nc$Best.n[1,])
    # barplot(table(nc$Best.n[1,]),
    #     xlab = "Number of Clusters", ylab = "Number of Criteria",
    #     main = "Number of Clusters Chosen by 26 Criteria")
