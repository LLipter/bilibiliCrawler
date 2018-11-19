library(DBI)
library(RMySQL)
Sys.setlocale(locale="UTF-8")
con <- dbConnect(MySQL(), host="localhost", dbname="bilibili", user="root", password="57575207")

# data preprocessing
bangumi.datasize <- 100
rawdata <- dbGetQuery(con, paste("SELECT * FROM bangumi WHERE ABS(view-view_calculated) < view*0.1 AND epno>10 ORDER BY view DESC LIMIT", bangumi.datasize))
bangumi.followview <- matrix(unlist(rawdata[c("follow","view_calculated")]), 
                            nrow=bangumi.datasize, 
                            ncol=2,
                            dimnames=list(rawdata$title, c("follow","view")))
bangumi.data <- scale(bangumi.followview)
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
number.cluster <- 7
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
cluster.result <- rect.hclust(bangumi.hc, k=number.cluster, border = 2:number.cluster+1)
par(op)
dev.off()

# kmeans
library(cluster)
number.cluster <- 5
bangumi.kmeans <- kmeans(bangumi.data, number.cluster, nstart=10)
png(file="assets/kmeans_clustering.png",width=3000, height=3000, res=600, pointsize=9)
op <- par(family='STXihei')
oo <- options(scipen=10)
clusplot(bangumi.followview, bangumi.kmeans$cluster, 
        color=TRUE, shade=TRUE, 
        s.x.2d = list(x=bangumi.followview, labs=rownames(bangumi.followview), var.dec=NA),
        labels=0, 
        lines=0, 
        main="K-means Clustering on 100 Most Popular Anime in Bilibili",
        sub="",
        col.txt="black",
        xlab="Subscriber",
        ylab="View",
        cex = 0.5)
x <- c(bangumi.followview["齐木楠雄的灾难（日播&精选版）",1]+1.3e6,
        bangumi.followview["Re：从零开始的异世界生活",1]+1e6,
        bangumi.followview["工作细胞",1]-0.4e6,
        bangumi.followview["紫罗兰永恒花园",1]+0.6e6,
        bangumi.followview["OVERLORD",1]-0.5e6,
        bangumi.followview["食戟之灵",1]-0.4e6,
        bangumi.followview["齐木楠雄的灾难 第二季",1]+0.9e6
        )
y <- c(bangumi.followview["齐木楠雄的灾难（日播&精选版）",2],
        bangumi.followview["Re：从零开始的异世界生活",2]-0.4e7,
        bangumi.followview["工作细胞",2]-0.4e7,
        bangumi.followview["紫罗兰永恒花园",2]-0.4e7,
        bangumi.followview["OVERLORD",2]+0.3e7,
        bangumi.followview["食戟之灵",2]+0.3e7,
        bangumi.followview["齐木楠雄的灾难 第二季",2]-0.3e7
        )
l <- c("齐木楠雄的灾难（日播&精选版）",
        "Re：从零开始的异世界生活",
        "工作细胞",
        "紫罗兰永恒花园",
        "OVERLORD",
        "食戟之灵",
        "齐木楠雄的灾难 第二季"
        )
text(x,y,labels=l,cex=0.8)
par(op)
options(oo)
dev.off()




# pams
number.cluster <-5
bangumi.diss <- daisy(bangumi.data)
bangumi.pamv <- pam(bangumi.diss, number.cluster, diss = TRUE)
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