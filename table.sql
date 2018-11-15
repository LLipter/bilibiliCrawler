CREATE DATABASE bilibili character set utf8;

USE bilibili

CREATE TABLE video(
    aid INT,
    status INT,
    title VARCHAR(200),
    pubdate DATETIME,
    mid INT,
    cid INT,
    tid INT,
    view INT,
    dannmaku INT,
    reply INT,
    favorite INT,
    coin INT,
    share INT,
    his_rank INT,
    support INT,
    dislike INT,
    copyright INT,
    PRIMARY KEY(aid)
)charset=utf8;
