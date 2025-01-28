package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const (
	driverName = "mysql"
	dbName     = "bot_factory"
)

type MySQL struct {
	db *sql.DB
}

func New(username, password, host string, port int) *MySQL {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", username, password, host, port, dbName)
	db := mustConnect(dataSource)
	return &MySQL{
		db: db,
	}
}

func mustConnect(dataSource string) *sql.DB {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		panic("couldn't connect to mysql driver: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("couldn't connect to mysql driver: " + err.Error())
	}
	return db
}
