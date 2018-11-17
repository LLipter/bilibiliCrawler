package db

import (
	"database/sql"
	"fmt"
	"github.com/LLipter/bilibiliCrawler/conf"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var (
	connPool *sql.DB // database connection poold
)

func init() {
	var err error
	connPool, err = sql.Open("mysql", conf.DBConnStr)
	if err != nil {
		fmt.Printf("cannot create database connection pool, %v\n", err)
		os.Exit(1)
	}
	err = connPool.Ping()
	if err != nil {
		fmt.Printf("cannot access database, %v\n", err)
		os.Exit(1)
	}
	connPool.SetMaxOpenConns(conf.DBconfig.MaxOpenConn)
	// https://github.com/go-sql-driver/mysql/issues/257
	/*
		Had the same problem:
		Didn't want to set SetMaxIdleConns to 0 since that would effectively disable connection pooling.
		What I did was changed SetConnMaxLifetime to a value that was less than the setting on MYSQL server
		that closed connections that were idle. MySQL was set to 10seconds before it closed an idle connection
		so I changed it to 9 seconds. Bug was fixed. The default value keeps idle connections as long as possible,
		so the connection pool thinks a particular connection is alive when in fact MySQL had already closed it.
		xThe connection pool then attempts to use the connection and you get the error.
	*/
	// that's why I choose this magic number....
	connPool.SetConnMaxLifetime(time.Second * 9)

}

func CloseDatabase() {
	if connPool != nil {
		connPool.Close()
	}
}
