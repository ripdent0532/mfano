package model

import (
	"database/sql"
	"time"

	logger "github.com/rip0532/mfano/lib/log"
)

type ProjectMapper struct {
	Id        int
	Name      string
	Created   time.Time
	Url       string
	CreatorId int64 `db:"creator_id"`
	Creator   string
	Group     string
	GroupId   int64 `db:"group_id"`
}

type ProjectSelectMap struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Created   time.Time `json:"created"`
	Url       string    `json:"url"`
	CreatorId int       `json:"creator_id" db:"creator_id"`
	Creator   string    `json:"creator"`
	Group     string    `json:"group"`
	GroupId   int64     `db:"group_id"`
}

func (mapper *ProjectMapper) Add() (sql.Result, error) {
	sql := "INSERT INTO mfano_project(name, created, url, creator_id, creator, `group`, group_id) values(:name, :created, :url, :creator_id, :creator, :group, :group_id)"
	result, err := DB.NamedExec(sql, &mapper)
	if err != nil {
		logger.Error.Printf("新增项目出错，错误信息：%v\n", err.Error())
	}
	return result, err
}

func (mapper *ProjectMapper) Select(groupId string, pageNum, pageSize int) ([]ProjectSelectMap, error) {
	var result []ProjectSelectMap
	s := struct {
		GroupId    string `db:"group_id"`
		PageOffset int    `db:"page_offset"`
		PageSize   int    `db:"page_size"`
	}{
		GroupId:    groupId,
		PageSize:   pageSize,
		PageOffset: pageNum * pageSize,
	}
	sql := `select * from mfano_project`
	if groupId != "" {
		sql = sql + ` where group_id = :group_id`
	}
	sql = sql + ` order by created desc limit :page_size offset :page_offset`
	rows, err := DB.NamedQuery(sql, s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		resultMap := ProjectSelectMap{}
		if err := rows.StructScan(&resultMap); err != nil {
			return nil, err
		}
		result = append(result, resultMap)
	}

	return result, nil
}

func (mapper *ProjectMapper) Count(groupId string) (int64, error) {
	sql := `select count(id) from mfano_project`
	if groupId != "" {
		sql = sql + ` where group_id = :group_id`
	}
	s := struct {
		GroupId string `db:"group_id"`
	}{
		GroupId: groupId,
	}
	rows, err := DB.NamedQuery(sql, &s)
	if err != nil {
		logger.Error.Printf("[DB] Count project error, error message: %v\n", err.Error())
		return 0, err
	}
	var count int64
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			logger.Error.Printf("[DB] Scan project error, error message: %v\n", err.Error())
			return 0, err
		}
	}
	return count, nil
}
