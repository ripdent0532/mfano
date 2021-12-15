package service

import (
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rip0532/mfano/lib/constant"
	"github.com/rip0532/mfano/lib/gravatar"
	"github.com/rip0532/mfano/lib/secret"
	"github.com/rip0532/mfano/middleware"
	"github.com/rip0532/mfano/model"
	"github.com/rip0532/mfano/serializer"
)

type UserListService struct {
	UserName string `form:"username"`
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
}

func (service *UserListService) Query() serializer.Response {
	mapper := &model.UserMapper{}
	result, err := mapper.Select(service.UserName, service.PageNum, service.PageSize)
	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "查询用户失败",
			Error: err.Error(),
		}
	}
	count, err := mapper.Count(service.UserName)
	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "查询用户失败",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 0,
		Data: result,
		Page: serializer.Page{Total: count, Size: service.PageSize, Num: service.PageNum},
	}
}

type UserLoginService struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Resource struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func (service *UserLoginService) Login(context *gin.Context) serializer.Response {
	userMapper := model.UserMapper{}
	user, err := userMapper.GetUserByName(service.UserName)
	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "用户登录失败",
			Error: err.Error(),
		}
	}
	if user.Password != "" && secret.CheckPassword(user.Password, service.Password) {
		session := sessions.Default(context)
		avatarURL := gravatar.New(user.Email).Size(200).Rating(gravatar.Pg).AvatarURL()
		session.Set("user", middleware.SessionUser{Name: user.Name, NickName: user.Nickname, Id: user.Id, AvatarURL: avatarURL, Email: user.Email})
		session.Save()
		resources := make([]Resource, 0, 100)
		if strings.EqualFold(user.Name, "admin") {
			// 返回管理员资源
			resources = append(resources, Resource{Url: "/user.html", Name: "用户"})
		}

		return serializer.Response{
			Code: 0,
			Data: gin.H{"resources": resources, "user": session.Get("user")},
		}
	}
	return serializer.Response{
		Code: 40001,
		Msg:  "用户名或密码错误",
	}
}

type UserRegisterService struct {
	Name     string `json:"username" binding:"required"`
	NickName string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (service *UserRegisterService) Register() serializer.Response {
	// 检查邮箱
	userMapper := model.UserMapper{}
	present, err := userMapper.EmailPresent(service.Email)
	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "用户注册失败",
			Error: err.Error(),
		}
	}
	if present {
		return serializer.Response{
			Code:  40001,
			Data:  nil,
			Msg:   "邮箱已注册",
			Error: "",
		}
	}

	// 检查用户名
	present, err = userMapper.NamePresent(service.Name)
	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "用户名已存在",
			Error: err.Error(),
		}
	}
	if present {
		return serializer.Response{
			Code:  40001,
			Msg:   "用户注册失败",
			Error: err.Error(),
		}
	}

	userMapper.Name = service.Name
	userMapper.NickName = service.NickName
	userMapper.Password = secret.EncodePassword(constant.DefaultPassword)
	userMapper.Email = service.Email
	_, err = userMapper.Add()

	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "注册失败",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 0,
		Msg:  "注册成功",
	}
}

type UserUpdateService struct {
	Name            string `json:"username" binding:"required"`
	NickName        string `json:"nickname" binding:"required"`
	Email           string `json:"email" binding:"required"`
	OriginPassword  string `json:"originPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`

	UserId int64
}

func (service *UserUpdateService) Update() serializer.Response {
	// 查询用户原始密码
	mapper := &model.UserMapper{}
	password, err := mapper.GetPassword(service.UserId)
	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "更新用户失败",
			Error: err.Error(),
		}
	}
	if service.OriginPassword != "" {
		if !secret.CheckPassword(password, service.OriginPassword) {
			return serializer.Response{
				Code: 40001,
				Msg:  "原密码错误，请重新输入",
			}
		} else {
			password = secret.EncodePassword(service.NewPassword)
		}
	}

	mapper.Id = service.UserId
	mapper.Name = service.Name
	mapper.NickName = service.NickName
	mapper.Password = password
	mapper.Email = service.Email

	if err := mapper.Update(); err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "更新用户信息错误",
			Error: err.Error(),
		}
	}
	return serializer.Response{
		Code: 0,
		Msg:  "更新用户信息成功",
	}
}

type UserAddService struct {
	Name     string `json:"username" binding:"required,min=3,max=10"`
	NickName string `json:"nickname" binding:"required,min=3,max=10"`
	Email    string `json:"email" binding:"required,email"`
}

func (service *UserAddService) Add() serializer.Response {
	mapper := model.UserMapper{}
	present, err := mapper.EmailPresent(service.Email)
	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "用户注册失败",
			Error: err.Error(),
		}
	}
	if present {
		return serializer.Response{
			Code: 40001,
			Msg:  "邮箱已注册",
		}
	}

	present, err = mapper.NamePresent(service.Name)
	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "用户注册失败",
			Error: err.Error(),
		}
	}
	if present {
		return serializer.Response{
			Code: 40001,
			Msg:  "用户名已注册",
		}
	}

	mapper.Name = service.Name
	mapper.NickName = service.NickName
	mapper.Password = secret.EncodePassword(constant.DefaultPassword)
	mapper.Email = service.Email
	_, err = mapper.Add()
	if err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "用户注册失败",
			Error: err.Error(),
		}
	}
	return serializer.Response{
		Code: 0,
		Msg:  "新增用户成功",
	}
}

type UserForbiddenService struct {
	Activity bool
	Id       int64
}

func (service *UserForbiddenService) forbidden() {
	mapper := &model.UserMapper{}
	mapper.UpdateActivity(service.Activity, service.Id)
}
