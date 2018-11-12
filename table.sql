CREATE DATABASE bilibili;

CREATE TABLE video(
    aid INT,
    status INT,
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
    cid INT,
    PRIMARY KEY(aid,cid),
    FOREIGN KEY(cid) REFERENCES video(aid)
)charset=utf8;