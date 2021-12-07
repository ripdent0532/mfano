package model

import (
	"fmt"

	"github.com/rip0532/mfano/lib/db"
	logger "github.com/rip0532/mfano/lib/log"
	"github.com/rip0532/mfano/lib/page"
	"github.com/rip0532/mfano/lib/secret"
)

type User struct {
	Id       int
	Name     string
	NickName string
	Password string
	Email    string
}

type PageUser struct {
	Users []*User
	Total int
	Num   int
	Size  int
}

func NewUser() *User {
	return &User{}
}

func (u *User) Add() error {
	conn := db.Open()
	stmt, err := conn.Prepare("INSERT INTO USER (NAME, NICKNAME, PASSWORD, EMAIL) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.Name, u.NickName, u.Password, u.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUsers(page page.Page) *PageUser {
	conn := db.Open()
	var queryStr string
	if u.Name != "" {
		queryStr = fmt.Sprintf("SELECT ID, NAME, NICKNAME, EMAIL FROM USER WHERE NAME='%s' ORDER BY ID ASC LIMIT %d OFFSET %d", u.Name, page.Size, page.Size*page.Num)
	} else {
		queryStr = fmt.Sprintf("SELECT ID, NAME, NICKNAME, EMAIL FROM USER ORDER BY ID ASC LIMIT %d OFFSET %d", page.Size, page.Size*page.Num)
	}
	rows, _ := conn.Query(queryStr)
	defer rows.Close()
	users := make([]*User, 0, page.Size)
	for rows.Next() {
		user := &User{}
		rows.Scan(&user.Id, &user.Name, &user.NickName, &user.Email)
		users = append(users, user)
	}
	var countQueryStr string
	if u.Name != "" {
		countQueryStr = fmt.Sprintf("SELECT COUNT(1) FROM USER WHERE NAME='%s'", u.Name)
	} else {
		countQueryStr = "SELECT COUNT(1) FROM USER"
	}
	var count int
	row := conn.QueryRow(countQueryStr)
	row.Scan(&count)
	return &PageUser{
		Users: users,
		Total: count,
		Num:   page.Num,
		Size:  page.Size,
	}
}

func (u *User) Get(name string) *User {
	conn := db.Open()
	queryStr := fmt.Sprintf("SELECT id, name, nickname, password, email FROM USER WHERE NAME = '%s'", name)
	row := conn.QueryRow(queryStr)
	user := &User{}
	row.Scan(&user.Id, &user.Name, &user.NickName, &user.Password, &user.Email)
	return user
}

func (u *User) ChangePassword(id int, password string) error {
	conn := db.Open()
	stmt, err := conn.Prepare("UPDATE USER SET PASSWORD = ? WHERE ID = ?")
	if err != nil {
		logger.Error.Println(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(secret.EncodePassword(password), id)
	return err
}

func (u *User) Update() error {
	conn := db.Open()
	stmt, err := conn.Prepare("UPDATE USER SET NAME = ?, NICKNAME = ?, PASSWORD = ?, EMAIL = ? WHERE ID = ?")
	defer stmt.Close()
	_, err = stmt.Exec(u.Name, u.NickName, u.Password, u.Email, u.Id)
	return err
}
