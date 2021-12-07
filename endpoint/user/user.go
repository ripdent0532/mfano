package user

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rip0532/mfano/lib/constant"
	"github.com/rip0532/mfano/lib/gravatar"
	logger "github.com/rip0532/mfano/lib/log"
	"github.com/rip0532/mfano/lib/page"
	"github.com/rip0532/mfano/lib/secret"
	"github.com/rip0532/mfano/middleware"
	"github.com/rip0532/mfano/model"
)

func Register(group *gin.RouterGroup) {
	group.GET("/user", userInfo)
	group.GET("/users", userList)
	group.POST("/user", addUser)
	group.POST("/user/change_password", changePassword)
	group.GET("/user/avatar_profile", avatarProfile)
	group.POST("/user/update", update)
}

func update(context *gin.Context) {
	var reqParam ChangeUserInfoReqParam
	if err := context.ShouldBind(&reqParam); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": 40001,
			"msg":  "操作失败，请检查参数",
			"err":  err.Error(),
		})
		return
	}
	currentUser := model.NewUser().Get(middleware.GetLoginUser(context).Name)
	if reqParam.Name != "" {
		currentUser.Name = reqParam.Name
	}
	if reqParam.NickName != "" {
		currentUser.NickName = reqParam.NickName
	}
	if reqParam.Email != "" {
		currentUser.Email = reqParam.Email
	}
	if reqParam.OriginPassword != "" {
		if !secret.CheckPassword(currentUser.Password, reqParam.OriginPassword) {
			context.JSON(http.StatusNotAcceptable, gin.H{
				"message": "原始密码错误",
			})
			return
		} else {
			currentUser.Password = secret.EncodePassword(reqParam.NewPassword)
		}
	}
	currentUser.Update()
	context.JSON(http.StatusOK, "success")
}

type Resource struct {
	Url  string
	Name string
}

func userInfo(context *gin.Context) {
	session := sessions.Default(context)
	user := session.Get("user")
	realUser := user.(middleware.SessionUser)
	resources := make([]Resource, 0, 100)
	if strings.EqualFold(realUser.Name, "admin") {
		// 返回管理员资源
		resources = append(resources, Resource{Url: "/user.html", Name: "用户"})
	}
	context.JSON(http.StatusOK, gin.H{
		"user":      realUser,
		"resources": resources,
	})
}

func userList(context *gin.Context) {
	page := page.Default()
	var reqParam QueryReqParam
	if err := context.ShouldBindQuery(&reqParam); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": 40001,
			"msg":  "操作失败，请检查参数",
			"err":  err.Error(),
		})
		return
	}

	page.Num = reqParam.PageNum
	user := model.NewUser()
	if reqParam.UserName != "" {
		user.Name = reqParam.UserName
	}
	result := user.GetUsers(page)

	context.JSON(http.StatusOK, gin.H{
		"data":  result.Users,
		"total": result.Total,
		"num":   result.Num,
		"size":  result.Size,
	})
}

func addUser(context *gin.Context) {
	var reqParam AddUserReqParam
	if err := context.ShouldBind(&reqParam); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": 40001,
			"msg":  "操作失败，请检查参数",
			"err":  err.Error(),
		})
		return
	}
	user := model.NewUser()
	user.Name = reqParam.Name
	user.NickName = reqParam.NickName
	user.Password = secret.EncodePassword(constant.DefaultPassword)
	user.Email = reqParam.Email
	if err := user.Add(); err != nil {
		logger.Error.Printf("新增用户出错：%s", err.Error())
		context.JSON(http.StatusBadRequest, nil)
	} else {
		context.JSON(http.StatusOK, nil)
	}
}

func changePassword(context *gin.Context) {
	loginUser := middleware.GetLoginUser(context)
	password := context.Request.PostFormValue("password")
	err := model.NewUser().ChangePassword(loginUser.Id, password)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"data": "更新密码失败"})
	} else {
		context.JSON(http.StatusOK, gin.H{"data": "OK"})
	}
}

func avatarProfile(context *gin.Context) {
	user := middleware.GetLoginUser(context)
	url := gravatar.New(user.Email).JSONURL()
	context.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}
