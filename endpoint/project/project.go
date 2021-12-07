package project

import (
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rip0532/mfano/lib"
	"github.com/rip0532/mfano/lib/archive"
	"github.com/rip0532/mfano/lib/constant"
	"github.com/rip0532/mfano/lib/db"
	logger "github.com/rip0532/mfano/lib/log"
	"github.com/rip0532/mfano/lib/page"
	"github.com/rip0532/mfano/middleware"
	"github.com/rip0532/mfano/model"
)

func Register(group *gin.RouterGroup) {
	group.POST("/project", add)
	group.GET("/projects", list)
}

// add 新增project
func add(context *gin.Context) {
	files, err := upload(context)
	zip := archive.NewZip()
	if err != nil {
		logger.Error.Printf("upload file error. %s\n", err)
	}
	for _, file := range files {

		err := zip.UnZip(constant.DstDir, file)
		if err != nil {
			logger.Error.Printf("upload file error. %s\n", err)
			context.JSON(http.StatusBadRequest, "文件处理失败")
			// 清理文件
			lib.RemoveFileOrFolder(constant.DstDir + "/" + strings.TrimSuffix(file, ".zip"))
			return
		}
	}
	insert(files, context)
	context.JSON(http.StatusOK, gin.H{"data": "OK"})
}

func upload(context *gin.Context) (files []string, err *error) {
	mediaType, params, _ := mime.ParseMediaType(context.Request.Header.Get("Content-Type"))
	if strings.HasPrefix(mediaType, "multipart/") {
		files = make([]string, 0, 10)
		multiReader := multipart.NewReader(context.Request.Body, params["boundary"])
		for {
			p, err := multiReader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			file, _ := os.OpenFile(constant.HomeDir+"/"+p.FileName(), os.O_WRONLY|os.O_CREATE, 0666)
			io.Copy(file, p)
			file.Close()
			files = append(files, p.FileName())
		}
	}
	return files, nil
}

func insert(files []string, context *gin.Context) {
	created := time.Now()
	loginUser := middleware.GetLoginUser(context)

	var addRepParam AddReqParam
	if err := context.ShouldBindQuery(&addRepParam); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": 40001,
			"msg":  "操作失败，请检查参数",
			"err":  err.Error(),
		})
		return
	}

	groupModel := model.NewGroup()
	groupId := model.NewGroup().GetGroupId(addRepParam.Group)
	if groupId == 0 {
		groupModel.Name = addRepParam.Group
		groupModel.Created = created
		groupModel.Creator = loginUser.NickName
		groupModel.CreatorId = loginUser.Id
		groupId, _ = groupModel.AddGroup()
	}

	for _, file := range files {
		project := model.NewProject()
		file = strings.TrimSuffix(file, ".zip")
		id, err := project.GetProjectId(file)
		if err != nil {
			logger.Error.Println(err)
			context.JSON(http.StatusBadRequest, err)
		}
		project.Name = file
		project.Creator = loginUser.NickName
		project.CreatorId = loginUser.Id
		project.Created = created
		project.GroupId = groupId
		project.Group = addRepParam.Group
		if id == 0 {
			project.AddProject()
		} else {
			project.Id = id
			project.UpdateProject()
		}
	}
}

type Project struct {
	Id      int
	Name    string
	Created string
	Url     string
	Creator string
	Group   string
}

type ProjectResult struct {
	Id      int
	Name    string
	Created time.Time
}

// list 查询列表
func list(context *gin.Context) {
	var param queryReqParam
	if err := context.ShouldBindQuery(&param); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": 40001,
			"msg":  "操作失败，请检查参数",
			"err":  err.Error(),
		})
		return
	}

	conn := db.Open()

	page := page.Default()
	page.Num = param.PageNum
	var queryStr string
	if param.GroupId != "" {
		queryStr = fmt.Sprintf("SELECT id, name, created, url, creator, `group` FROM project where group_id = %s order by created desc limit %d offset %d", param.GroupId, page.Size, page.Num*page.Size)
	} else {
		queryStr = fmt.Sprintf("SELECT id, name, created, url, creator, `group` FROM project order by created desc limit %d offset %d", page.Size, page.Num*page.Size)
	}
	rows, err := conn.Query(queryStr)
	println(err)
	defer rows.Close()
	result := make([]*Project, 0, 100)
	for rows.Next() {
		var id int
		var name string
		var created time.Time
		var url string
		var creator string
		var group string
		rows.Scan(&id, &name, &created, &url, &creator, &group)
		p := &Project{
			Id:      id,
			Name:    name,
			Created: created.Format("2006-01-02 15:04:05"),
			Url:     url,
			Creator: creator,
			Group:   group,
		}
		result = append(result, p)
	}
	// 统计总数
	var count int
	var countQueryStr string
	if param.GroupId != "" {
		countQueryStr = fmt.Sprintf("SELECT COUNT(ID) FROM PROJECT where group_id = %s", param.GroupId)
	} else {
		countQueryStr = "SELECT COUNT(ID) FROM PROJECT"
	}
	row := conn.QueryRow(countQueryStr)
	row.Scan(&count)
	context.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": count,
		"num":   page.Num,
		"size":  page.Size,
	})
}
