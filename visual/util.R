# Within cluster sum of squares by cluster
wssplot <- function(data, nc=15, ns=10){
    wss <- NULL
    for(i in 1:nc){
      wss[i] <- sum(kmeans(data, centers=i, nstart=ns)$withinss)
    }
    plot(1:nc, wss, type = "b",xlab = "Number of Clusters",
        ylab = "Within groups sum of squares")
}

# number of indices ~ number of clusters
nincplot <- function(data, method="kmeans"){
    op <- par(no.readonly=TRUE)
    res <- NbClust(data, method=method)
    par(op)
    graph <- table(res$Best.nc[1,])
    barplot(graph,
        xlab = "Number of Clusters", ylab = "Number of Indices")
}