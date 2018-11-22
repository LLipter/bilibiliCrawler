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
    pages INT,
    PRIMARY KEY(aid)
)charset=utf8;

CREATE INDEX status_idx ON video(status);
CREATE INDEX pubdate_idx ON video(pubdate);
CREATE INDEX dannmaku_idx ON video(dannmaku);
CREATE INDEX view_idx ON video(view);
CREATE INDEX pages_idx ON video(pages);
CREATE INDEX tid_idx ON video(tid);

CREATE TABLE online(
    ts         DATETIME,
	douga      INT, -- 动画
	anime      INT, -- 番剧
	guochuang  INT, -- 国创
	music      INT, -- 音乐
	dance      INT, -- 舞蹈
	game       INT, -- 游戏
	technology INT, -- 科技
	life       INT, -- 生活
	kichiku    INT, -- 鬼畜
	fashion    INT, -- 时尚
	ad         INT, -- 广告
	ent        INT, -- 娱乐
	cinephile  INT, -- 影视
	cinema     INT, -- 放映厅
	tv         INT, -- 电视剧
	movie      INT, -- 电影
    
    allcount   INT, -- 最新投稿
	webonline  INT, -- 在线人数
	playonline INT, 
    PRIMARY KEY(ts)
)charset=utf8;

CREATE TABLE bangumi(
    sid INT,
    title VARCHAR(200),
    pubdate DATETIME,
    epno INT,
    score DOUBLE(3,2),
    follow INT,
    view INT,
    view_calculated INT,
    media_id INT,
    PRIMARY KEY(sid)
)charset=utf8;

CREATE TABLE episode(
    sid   INT,
    idx   INT,
    aid   INT,
    view  INT,
	cid   INT,
	epid  INT,
    PRIMARY KEY(sid,idx)
)charset=utf8;
