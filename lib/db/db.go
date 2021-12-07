package db

import (
	"database/sql"
	"sync"

	"github.com/rip0532/mfano/lib/constant"
	logger "github.com/rip0532/mfano/lib/log"
)

var (
	instance *sql.DB
	once     sync.Once
)

func Open() *sql.DB {
	once.Do(func() {
		if db, err := sql.Open("sqlite3", constant.Db_Dir+"/mfano.sqlite"); err != nil {
			logger.Error.Println(err.Error())
			panic(err)
		} else {
			maxOpenConns := 20
			maxIdleConns := 5
			db.SetMaxOpenConns(maxOpenConns)
			db.SetMaxIdleConns(maxIdleConns)
			instance = db
			logger.Info.Printf("初始化数据库连接, MaxOpenConns: %d MaxInleConns: %d\n", maxOpenConns, maxIdleConns)
		}
	})
	return instance
}
