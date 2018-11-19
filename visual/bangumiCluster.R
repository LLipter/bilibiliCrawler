library(DBI)
library(RMySQL)
Sys.setlocale(locale="UTF-8")
con <- dbConnect(MySQL(), host="localhost", dbname="bilibili", user="root", password="57575207")

# data preprocessing
rawdata <- dbGetQuery(con, "SELECT * FROM bangumi WHERE ABS(view-view_calculated) < view*0.1 AND epno>10 ORDER BY view DESC LIMIT 200")
bangumi.data <- rawdata[c("follow", "view_calculated")]
bangumi.data <- scale(bangumi.data)
row.names(bangumi.data) <- rawdata$title
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
bangumi.data <- cbind(bangumi.data, viewdata)

# Hierarchical Clustering
bangumi.dist <- dist(bangumi.data[1:50,])
bangumi.hc <- hclust(bangumi.dist)
number_cluster <- 7
png(file="assets/hierarchical_clustering.png",width=3000, height=3000, res=600, pointsize=9)
op <- par(mai=c(0.1,0.7,0.1,0.1),lwd=0.7,font.lab=2)
plot(bangumi.hc, 
    hang = -1, 
    cex = .5, 
    ylab = "Hierarchical Clustering on 50 Most Popular Anime in Bilibili",
    family='STXihei',
    xlab="",
    main="",
    sub="",
    axes=TRUE,
    )
cluster.result <- rect.hclust(bangumi.hc, k=number_cluster, border = 2:number_cluster+1)
par(op)
dev.off()


# kmeans
number_cluster <- 5
bangumi.kmeans <- kmeans(bangumi.data, number_cluster, nstart=10)
library(cluster)
png(file="assets/kmeans_clustering.png",width=3000, height=3000, res=600, pointsize=9)
op <- par(family='STXihei')
clusplot(bangumi.data[,1:2], bangumi.kmeans$cluster, 
        color=TRUE, shade=FALSE, 
        labels=0, 
        lines=0, 
        main="K-means Clustering on 200 Most Popular Anime in Bilibili",
        col.txt="black",
        cex = 0.5)
clusplot(bangumi.kmeans$cluster, 
        color=TRUE, shade=FALSE, 
        labels=0, 
        lines=0, 
        main="K-means Clustering on 200 Most Popular Anime in Bilibili",
        col.txt="black",
        cex = 0.5)
par(op)
dev.off()


# pams
bangumi.diss <- daisy(bangumi.data)
bangumi.pamv <- pam(bangumi.diss, number_cluster, diss = TRUE)
op <- par(family='STXihei')
library(cluster)
clusplot(bangumi.pamv, 
        shade=FALSE, color=TRUE, 
        lines=0,
        labels=0,
        main="PAM Clustering on 200 Most Popular Anime in Bilibili",
        col.txt="black",
        cex = 0.5)
par(op)

dbDisconnect(con)






x <- rbind(cbind(rnorm(10,0,0.5), rnorm(10,0,0.5)),
           cbind(rnorm(15,5,0.5), rnorm(15,5,0.5)))
pamx <- pam(x, 2)
pamx # Medoids: '7' and '25' ...
summary(pamx)
plot(pamx)