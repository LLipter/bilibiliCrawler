# load liabrary
library(DBI)
library(RMySQL)
library(cluster)

# connect to database
con <- dbConnect(MySQL(), host="localhost", dbname="bilibili", user="(づ｡◕‿‿◕｡)づ", password="123456")

# setting parameters
bangumi.datasize <- 100
number.hc.datasize <- 50
number.hc.cluster <- 7
number.kmeans.datasize <- 100
number.kmeans.cluster <- 5
number.pam.datasize <- 100
number.pam.cluster <- 5
Sys.setlocale(locale="UTF-8")

# data preprocessing
bangumi.rawdata <- dbGetQuery(con, paste("SELECT * FROM bangumi WHERE ABS(view-view_calculated) < view*0.1 AND epno>10 ORDER BY view DESC LIMIT", bangumi.datasize))
bangumi.followview <- matrix(unlist(bangumi.rawdata[c("follow","view_calculated")]), 
                            nrow=bangumi.datasize, 
                            ncol=2,
                            dimnames=list(bangumi.rawdata$title, c("follow","view")))
bangumi.data <- scale(bangumi.followview)
viewdata <- NULL  
for(sid in bangumi.rawdata$sid){
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
bangumi.dist <- dist(bangumi.data[1:number.hc.datasize,])
bangumi.hc <- hclust(bangumi.dist)
png(file="assets/hierarchical_clustering.png", width=3000, height=3000, res=600, pointsize=9)
op <- par(mai=c(0.1,0.7,0.1,0.1),lwd=0.7)
plot(bangumi.hc, 
    hang = -1, 
    cex = .5, 
    ylab = paste("Hierarchical Clustering on", number.hc.datasize, "Most Popular Anime in Bilibili"),
    family='STXihei',
    xlab="",
    main="",
    sub="",
    axes=TRUE)
cluster.result <- rect.hclust(bangumi.hc, k=number.hc.cluster, border = 2:number.hc.cluster+1)
par(op)
dev.off()

# kmeans
bangumi.kmeans <- kmeans(bangumi.data[1:number.kmeans.datasize,], number.kmeans.cluster, nstart=50)
png(file="assets/kmeans_clustering.png",width=3000, height=3000, res=600, pointsize=9)
op <- par(family='STXihei')
oo <- options(scipen=10)
clusplot(bangumi.followview[1:number.kmeans.datasize,], bangumi.kmeans$cluster, 
        color=TRUE, shade=FALSE, 
        s.x.2d = list(x=bangumi.followview[1:number.kmeans.datasize,], labs=rownames(bangumi.followview)[1:number.kmeans.datasize], var.dec=NA),
        labels=0, 
        lines=0, 
        main=paste("Kmeans Clustering on", number.kmeans.datasize, "Most Popular Anime in Bilibili"),
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
y <- c(bangumi.followview["齐木楠雄的灾难（日播&精选版）",2]+0.2e7,
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
bangumi.diss <- daisy(bangumi.data[1:number.pam.datasize,])
bangumi.pamv <- pam(bangumi.diss, number.pam.cluster, diss = TRUE)
png(file="assets/pam_clustering.png",width=3000, height=3000, res=600, pointsize=9)
op <- par(family='STXihei')
oo <- options(scipen=10)
clusplot(bangumi.followview[1:number.pam.datasize,], bangumi.pamv$clustering, 
        shade=FALSE, color=TRUE, 
        s.x.2d = list(x=bangumi.followview[1:number.pam.datasize,], labs=rownames(bangumi.followview)[1:number.pam.datasize], var.dec=NA),
        lines=0,
        labels=0,
        main=paste("PAM Clustering on", number.pam.datasize, "Most Popular Anime in Bilibili"),
        col.txt="black",
        sub="",
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
y <- c(bangumi.followview["齐木楠雄的灾难（日播&精选版）",2]-0.2e7,
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

# silhouette plot of kmeans
kmeans.dist <- dist(bangumi.data[1:number.kmeans.datasize,])
kmeans.si <- silhouette(bangumi.kmeans$cluster, dist=kmeans.dist)
png(file="assets/kmeans_silhouette.png",width=5000, height=9000, res=600)
plot(kmeans.si,
    col=terrain.colors(number.kmeans.datasize),
    main=paste("Silhouette Plot of Kmeans Clustering on", number.kmeans.datasize, "Most Popular Anime in Bilibili"),
    do.clus.stat=TRUE,
    do.col.sort=FALSE,
    do.n.k=FALSE,
    )
dev.off()

# silhouette plot of PAM
pam.si <- silhouette(bangumi.pamv)
png(file="assets/pam_silhouette.png",width=5000, height=9000, res=600)
plot(pam.si,
    col=terrain.colors(number.kmeans.datasize),
    main=paste("Silhouette Plot of PAM Clustering on", number.kmeans.datasize, "Most Popular Anime in Bilibili"),
    do.clus.stat=TRUE,
    do.col.sort=FALSE,
    do.n.k=FALSE,
    )
dev.off()

# close database
dbDisconnect(con)
