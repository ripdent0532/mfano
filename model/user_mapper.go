package model

import (
	"database/sql"

	logger "github.com/rip0532/mfano/lib/log"
)

type UserMapper struct {
	Id       int64
	Name     string
	NickName string `db:"nickname"`
	Password string
	Email    string
}

type UserSelectMap struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	NickName string `json:"nickname"`
	Email    string `json:"email"`
}

func (mapper *UserMapper) Select(name string, pageNum, pageSize int) ([]UserSelectMap, error) {
	var result []UserSelectMap
	s := struct {
		Name       string
		PageSize   int `db:"page_size"`
		PageOffset int `db:"page_offset"`
	}{
		Name:       name,
		PageSize:   pageSize,
		PageOffset: pageNum * pageSize,
	}
	sql := `select id, name, nickname, email from mfano_user`
	if name != "" {
		sql = sql + " where name = :name"
	}
	sql = sql + " order by id asc limit :page_size offset :page_offset"
	rows, err := DB.NamedQuery(sql, &s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		resultMap := UserSelectMap{}
		if err := rows.StructScan(&resultMap); err != nil {
			return nil, err
		}
		result = append(result, resultMap)
	}
	return result, nil
}

func (mapper *UserMapper) Count(name string) (int64, error) {
	s := struct {
		Name string
	}{Name: name}
	sql := `select count(1) from mfano_user`
	if name != "" {
		sql = sql + " where name = :name"
	}
	rows, err := DB.NamedQuery(sql, &s)
	if err != nil {
		logger.Error.Printf("[DB] Count user error, error info: %v\n", err.Error())
		return 0, err
	}
	defer rows.Close()
	var count int64
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			logger.Error.Printf("[DB] Count user error, error info: %v\n", err.Error())
			return 0, err
		}
	}
	return count, nil
}

type GetUserMap struct {
	Id       int64
	Name     string
	Nickname string
	Password string
	Email    string
}

func (mapper *UserMapper) GetUserByName(name string) (*GetUserMap, error) {
	result := &GetUserMap{}
	sql := `select id, name, nickname, password, email from mfano_user where name = ?`
	if err := DB.Get(result, sql, name); err != nil {
		logger.Error.Printf("[DB] Query user error, error info: %v\n", err)
		return nil, err
	}
	return result, nil
}

func (mapper *UserMapper) EmailPresent(email string) (bool, error) {
	sql := `select count(id) from mfano_user where email = ?`
	var count int
	if err := DB.Get(&count, sql, email); err != nil {
		logger.Error.Printf("[DB] Query user by email error, error info: %v\n", err)
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (mapper *UserMapper) NamePresent(name string) (bool, error) {
	sql := `select count(id) from mfano_user where name = ?`
	var count int
	if err := DB.Get(&count, sql, name); err != nil {
		logger.Error.Printf("[DB] Query user by name error, error info: %v\n", err)
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (mapper *UserMapper) Add() (sql.Result, error) {
	sql := `INSERT INTO MFANO_USER (NAME, NICKNAME, PASSWORD, EMAIL) VALUES (:name, :nickname, :password, :email)`
	sqlResult, err := DB.NamedExec(sql, mapper)
	if err != nil {
		logger.Error.Printf("[DB] Insert user error, error info: %v\n", err)
		return nil, err
	}
	return sqlResult, nil
}

func (mapper *UserMapper) GetPassword(id int64) (string, error) {
	var password string
	sql := `select password from mfano_user where id = ?`
	if err := DB.Get(&password, sql, id); err != nil {
		logger.Error.Printf("[DB] Get user password error, error info: %v\n", err)
		return "", err
	}
	return password, nil
}

func (mapper *UserMapper) Update() error {
	sql := `update mfano_user set name = :name, nickname = :nickname, password = :password, email = :email where id = :id`
	_, err := DB.NamedExec(sql, mapper)
	if err != nil {
		logger.Error.Printf("[DB] Update user error, error info: %v\n", err)
		return err
	}
	return nil
}

func (mapper *UserMapper) UpdateActivity(activity bool, id int64) error {
	sql := `update mfano_user set activity = ? where id = ?`
	_, err := DB.Exec(sql, activity, id)
	if err != nil {
		return err
	}
	return nil
}
