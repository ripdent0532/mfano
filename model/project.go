package model

import (
	"fmt"
	"time"

	"github.com/rip0532/mfano/lib/constant"
	"github.com/rip0532/mfano/lib/db"
	logger "github.com/rip0532/mfano/lib/log"
)

type Project struct {
	Id        int
	Name      string
	Created   time.Time
	Url       string
	CreatorId int
	Creator   string
	Group     string
	GroupId   int64
}

func NewProject() *Project {
	return &Project{}
}

func (project *Project) AddProject() {
	conn := db.Open()
	stmt, err := conn.Prepare("INSERT INTO project(name, created, url, creator_id, creator, `group`, group_id) values(?,?,?,?,?,?,?)")
	if err != nil {
		logger.Error.Println(err)
	}
	defer stmt.Close()
	stmt.Exec(project.Name, project.Created, constant.StaticHost+"/"+project.Name, project.CreatorId, project.Creator, project.Group, project.GroupId)
}

func (project *Project) UpdateProject() error {
	conn := db.Open()
	stmt, err := conn.Prepare("UPDATE project set created=?, creator_id=?, creator=?, group_id=?, `group`=? where id=?")
	if err != nil {
		logger.Error.Println(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(project.Created, project.CreatorId, project.Creator, project.GroupId, project.Group, project.Id)
	return nil
}

func (project *Project) GetProjectId(name string) (int, error) {
	conn := db.Open()
	row := conn.QueryRow(fmt.Sprintf("SELECT id FROM project where name='%s'", name))
	if row.Err() != nil {
		logger.Error.Println(row.Err())
		return 0, row.Err()
	}
	var id int
	row.Scan(&id)
	return id, nil
}
