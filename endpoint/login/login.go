package login

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rip0532/mfano/lib/gravatar"
	"github.com/rip0532/mfano/lib/secret"
	"github.com/rip0532/mfano/middleware"
	"github.com/rip0532/mfano/model"
)

func Register(e *gin.Engine, group *gin.RouterGroup) {
	group.GET("/logout", logout)
	e.POST("/login", login)
}

type Resource struct {
	Url  string
	Name string
}

// login 登录
func login(context *gin.Context) {
	var loginForm Form
	if err := context.ShouldBindJSON(&loginForm); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": 40001,
			"msg":  "登录失败，请检查参数",
			"err":  err.Error(),
		})
		return
	}

	username := loginForm.UserName
	loginPassword := loginForm.Password
	session := sessions.Default(context)
	user := model.NewUser().Get(username)
	match := secret.CheckPassword(user.Password, loginPassword)
	if match {
		avatarURL := gravatar.New(user.Email).Size(200).Rating(gravatar.Pg).AvatarURL()
		session.Set("user", middleware.SessionUser{Name: username, NickName: user.NickName, Id: user.Id, AvatarURL: avatarURL, Email: user.Email})
		session.Save()
		resources := make([]Resource, 0, 100)
		if strings.EqualFold(username, "admin") {
			// 返回管理员资源
			resources = append(resources, Resource{Url: "/user.html", Name: "用户"})
		}
		context.JSON(http.StatusOK, gin.H{"resources": resources, "user": session.Get("user")})
	} else {
		context.JSON(http.StatusBadRequest, "用户或密码错误")
	}

}

// logout 登出
func logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
}
