# Abstract

In this paper, I analyzed data collected from Bilibili, which is the most popular online entertainment website in China. By visualizing data and running clustering algorithm, some interesting results pop up. I use `Golang` to implement a crawler with high concurrency supported to collect data, and `R` takes over the analyzing and visualizing part if it.

# Introduction

An era of rapid economic growth in China with the boom of internet gave rise to a new generation - Generation Z, people who born between 1990 and 2009. Like no other generation before, they grow up in affluent times and are internet-savvy consumers and trend-setters. They are transforming online entertainment in China. Born to address their demand, Bilibili is the welcoming home of diverse culture and interests and represents the iconic brand for online entertainment serving young generations in China. Started as a content community inspired by anime, comics and games, Bilibili has evolved into a full-spectrum online entertainment world with 72 million and growing monthly-active-users, over 80% of whom are Generation z. Bilibili created an immersive entertainment experience for Generation Z, and built highly sticky and engaged communities. Their pioneered community feature - bullet chatting - transforms viewing experience by allowing audience to share thoughts and feelings real-time with others viewing the same video. They also created a unique membership exam, leading to users' strong sense of belonging. As a result, their users spend approximately 76 minutes daily on Bilibili. Their dynamic communities fuel on ever-growing supply of creative professional user generated content, and content creators earn respect and rewards from our users, forming a self-reinforcing virtuous cycle. By analyzing the data collected from it, I believe it's a good way to learn about the Generation Z.

# Part 1: Clustering Analyzing of Anime

### Data Collection

`bilibiliCrawler -b[d]`

Using the above command, over 2000 anime data will be collect from Bilibili. 

In this section, only the following field will be used.

| Field | Explanation | 
| --- | --- |
| title | name/identifier of the anime | 
| subscriber | the number of subscriber of the anime | 
| view | the number of times that the anime is played (sum of all episodes) |
| view1 | the number of times that the first episode of the anime is played |
| view3 | the number of times that the third episode of the anime is played |
| view1n | the number of times that the last episode of the anime is played |
| view3n | the number of times that the third last episode of the anime is played |

The view of first episode of an anime may just depends on the publicity and how popular of this intellectual property. So it may just demonstrate people's  expectation and be irrelevant to the actual performance of it. However, the view of third episode may largely relative with the actual quality of this anime. In this sense, by comparing the difference between `view1` and `view3`, I can find whether people is satisfied with it. The view of the third last episode may related with the equality of the entirely anime, and the view of the last episode may related with the quality of ending of this anime. I choose this four field in order to give more dimension information, so that my analysis will be more comprehensive and accurate.

### Data Preprocessing

I chose a subset of all anime to conduct this analysis, namely, 200 most popular anime were used as dataset after removing illegal data. Here I decided to use the number of times that the anime is played as criteria to determine which one is more popular, since people is generally more familiar with those anime, hence, the result will be more intuitive. 

Also, I use normalization to eliminate the influence caused by the order of magnitudes of different field. 

To be more specific, `view` and `subscriber` are normalized across **all anime** analyzed. The following codes does the job.

~~~R
# bangumi is anime in Japanese. 
# I use it just because it appears more interesting to me.
bangumi.data <- scale(bangumi.data)
~~~

`view1`, `view2`, `view1n` and `view3n` are a little different. They are normalized across **all episode** of a specific anime. The following codes does the job.

~~~R
epdata <- scale(epdata)
~~~

### Hierarchical Clustering

Hierarchical cluster analysis using a set of dissimilarities for the n objects being clustered. In this case, a dissimilarity matrix based on euclidean distance is computed. Initially, each anime is assigned to its own cluster and then the algorithm proceeds iteratively, at each stage joining the two most similar clusters, continuing until there is just a single cluster. At each stage distances between clusters are recomputed by the Lance–Williams dissimilarity update formula according to the particular clustering method being used.

~~~R
bangumi.dist <- dist(bangumi.data[1:50,])
bangumi.hc <- hclust(bangumi.dist)
~~~

![](visual/assets/hierarchical_clustering.png)

The above figure demonstrate the clustering result on 50 most popular anime. Here I choose a dataset of size 50 instead of 200 for better clearity.

Each colorful rectangle box in the figure identify a cluster after I divided them into 5 clusters based on the height of this hierarchical tree.

# Reference

1. [bilibili上市宣传视频](www.bilibili/video/av21322566)