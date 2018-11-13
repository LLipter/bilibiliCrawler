CREATE DATABASE bilibili;

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

CREATE TABLE chatid(
    aid INT,
    pageno INT,
    cid INT,
    duration INT,
    subtitle VARCHAR(200),
    PRIMARY KEY(aid,pageno),
    FOREIGN KEY(aid) REFERENCES video(aid)
)charset=utf8;