package service

import (
	"strings"
	"time"

	"github.com/rip0532/mfano/lib/archive"
	"github.com/rip0532/mfano/lib/constant"

	"github.com/rip0532/mfano/middleware"
	"github.com/rip0532/mfano/model"
	"github.com/rip0532/mfano/serializer"
)

type ProjectQueryService struct {
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size" page:"size=10"`
	GroupId  string `form:"group_id"`
}

func (service *ProjectQueryService) Query() serializer.Response {
	mapper := model.ProjectMapper{}
	result, err := mapper.Select(service.GroupId, service.PageNum, service.PageSize)
	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "查询项目失败",
			Error: err.Error(),
		}
	}
	count, err := mapper.Count(service.GroupId)
	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "查询项目失败",
			Error: err.Error(),
		}
	}
	return serializer.Response{
		Code: 0,
		Data: result,
		Page: serializer.Page{Total: count, Num: service.PageNum, Size: service.PageSize},
	}
}

type ProjectAddService struct {
	Group string `form:"group" binding:"required"`
	Files []archive.UploadFile
	User  middleware.SessionUser
}

func (service *ProjectAddService) Add() serializer.Response {
	groupMapper := model.GroupMapper{}
	result := groupMapper.Get(service.Group)
	created := time.Now()
	var groupId int64
	if result == nil {
		// 新增group
		groupMapper.Name = service.Group
		groupMapper.Created = created
		groupMapper.Creator = service.User.NickName
		groupMapper.CreatorId = service.User.Id

		addResult, err := groupMapper.Add()
		if err != nil {
			return serializer.Response{Code: 40001, Msg: "新增项目失败", Error: err.Error()}
		}
		groupId, _ = addResult.LastInsertId()
	} else {
		groupId = result.Id
	}

	// project url 格式：host/hash
	for _, file := range service.Files {
		fileName := strings.TrimSuffix(file.File, ".zip")
		projectMapper := model.ProjectMapper{}
		projectMapper.Name = fileName
		projectMapper.Creator = service.User.NickName
		projectMapper.CreatorId = service.User.Id
		projectMapper.Group = service.Group
		projectMapper.GroupId = groupId
		projectMapper.Url = constant.StaticHost + "/" + file.Timestamp + "/" + fileName
		projectMapper.Created = created

		_, err := projectMapper.Add()
		if err != nil {
			return serializer.Response{
				Code:  40001,
				Msg:   "新增项目失败",
				Error: err.Error(),
			}
		}
	}

	return serializer.Response{
		Code: 0,
		Msg:  "新增项目成功",
	}
}
