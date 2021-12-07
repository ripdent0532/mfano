package model

import (
	"fmt"
	"math"
	"time"

	"github.com/rip0532/mfano/lib/db"
)

type Group struct {
	Id        int
	Name      string
	Created   time.Time
	Creator   string
	CreatorId int
}

func NewGroup() *Group {
	return &Group{}
}

func (group *Group) List() []Group {
	conn := db.Open()
	rows, _ := conn.Query(fmt.Sprintf("SELECT * FROM `GROUP`"))
	results := make([]Group, 0, math.MaxInt8)
	defer rows.Close()
	for rows.Next() {
		result := Group{}
		rows.Scan(&result.Id, &result.Name, &result.Created, &result.Creator, &result.CreatorId)
		results = append(results, result)
	}
	return results
}

func (group *Group) GetGroupId(name string) int64 {
	conn := db.Open()
	row := conn.QueryRow(fmt.Sprintf("SELECT ID FROM `GROUP` WHERE NAME = '%s'", name))
	if row.Err() != nil {
		panic(row.Err())
	} else {
		var id int64
		row.Scan(&id)
		return id
	}
}

func (group *Group) AddGroup() (int64, error) {
	conn := db.Open()
	stmt, _ := conn.Prepare("INSERT INTO `GROUP` (NAME, CREATED, CREATOR, CREATOR_ID) VALUES (?, ?, ?, ?)")
	defer stmt.Close()
	result, _ := stmt.Exec(group.Name, group.Created, group.Creator, group.CreatorId)
	return result.LastInsertId()
}
