package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	logger "github.com/rip0532/mfano/lib/log"
)

var DB *sqlx.DB

func InitializeDatabase() (err error) {
	dsn := "root:root@tcp(localhost:3306)/mfano?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.Error.Printf("连接数据库失败，错误信息：%v\n", err)
		return err
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
	return nil
}
