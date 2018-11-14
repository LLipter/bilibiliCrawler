-- use `set names utf8` to ensure all charset is correct
CREATE DATABASE bilibili character set utf8;

USE bilibili

CREATE TABLE video(
    aid INT,
    status INT,
    title VARCHAR(200),
    pubdate DATETIME,
    mid INT,
    duration INT,
    view INT,
    dannmaku INT,
    reply INT,
    favorite INT,
    coin INT,
    share INT,
    now_rank INT,
    his_rank INT,
    support INT,
    dislike INT,
    no_reprint INT,
    copyright INT,
    PRIMARY KEY(aid)
)charset=utf8;

CREATE TABLE pages(
    aid INT,
    pageno INT,
    cid INT,
    duration INT,
    subtitle VARCHAR(200),
    PRIMARY KEY(aid,pageno)
)charset=utf8;
