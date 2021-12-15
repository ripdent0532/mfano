package model

import (
	"database/sql"
	"time"

	logger "github.com/rip0532/mfano/lib/log"
)

type GroupMapper struct {
	Id        int64
	Name      string
	Created   time.Time
	Creator   string
	CreatorId int64 `db:"creator_id"`
}

type GroupGetMap struct {
	Id        int64
	Name      string
	Created   time.Time
	Creator   string
	CreatorId int `db:"creator_id"`
}

func (mapper *GroupMapper) Get(groupName string) *GroupGetMap {
	result := &GroupGetMap{}
	sql := `select id, name, created, creator, creator_id from mfano_group where name = ?`
	if err := DB.Get(result, sql, groupName); err != nil {
		return nil
	}
	return result
}

func (mapper *GroupMapper) Add() (sql.Result, error) {
	sql := `INSERT INTO mfano_group (NAME, CREATED, CREATOR, CREATOR_ID) VALUES (:name, :created, :creator, :creator_id)`
	result, err := DB.NamedExec(sql, mapper)
	if err != nil {
		logger.Error.Printf("新增Group出错，错误信息: %v\n", err.Error())
		return nil, err
	}
	return result, nil

}
