library(DBI)
library(RMySQL)
Sys.setlocale(locale="UTF-8")
con <- dbConnect(MySQL(), host="localhost", dbname="bilibili", user="root", password="57575207")

# data preprocessing
rawdata <- dbGetQuery(con, "SELECT * FROM bangumi WHERE ABS(view-view_calculated) < view*0.1 AND epno>10 ORDER BY view DESC LIMIT 100")
bangumi_data <- rawdata[c("follow", "view_calculated")]
bangumi_data <- scale(bangumi_data)
row.names(bangumi_data) <- rawdata$title
viewdata <- NULL  
for(sid in rawdata$sid){
    sqlstr <- paste("SELECT view FROM episode WHERE sid =", sid, "ORDER BY idx")
    epdata <- dbGetQuery(con, sqlstr)
    epdata <- epdata$view
    epdata <- scale(epdata)
    len <- length(epdata)
    row <- c(epdata[1],epdata[3],epdata[len],epdata[len-3])
    viewdata <- rbind(viewdata, row)
}
colnames(viewdata) <- c("view1","view3","view1n","view3n") 
bangumi_data <- cbind(bangumi_data, viewdata)

# Hierarchical Clustering
bangumi_dist <- dist(bangumi_data)
bangumi_hc <- hclust(bangumi_dist, method = "average")
plot(bangumi_hc, 
    hang = -1, 
    cex = .5, 
    main = "Hierarchical Clustering on 100 Most Popular Anime in Bilibili",
    family='STXihei',
    xlab="",
    ylab=NULL,
    sub="",
    axes = FALSE)
cluster_result <- rect.hclust(bangumi_hc, k=10, border = "red")


# kmeans
number_cluster <- 7
bangumi_kmeans <- kmeans(bangumi_data, number_cluster, nstart=10)
plot(bangumi_data[,c("follow", "view_calculated")],
    col=bangumi_kmeans$cluster)
points(bangumi_kmeans$centers[,c("follow", "view_calculated")], col = 1:number_cluster , pch = 8, cex=2);


dbDisconnect(con)
