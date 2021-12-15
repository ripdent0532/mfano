package v1

import (
	"net/http"

	"github.com/rip0532/mfano/lib"

	"github.com/rip0532/mfano/middleware"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
	"github.com/rip0532/mfano/serializer"
	"github.com/rip0532/mfano/service"
)

//func UserRouterRegister(router *gin.RouterGroup) {
//	router.POST("/user/login", userLoginHandler)
//	router.POST("/user/register", userRegisterHandler)
//}

func UserLoginHandler(ctx *gin.Context) {
	var s service.UserLoginService
	if err := ctx.ShouldBindJSON(&s); err != nil {
		ctx.JSON(http.StatusOK, serializer.Response{
			Code:  40001,
			Msg:   "操作失败，请检查参数",
			Error: err.Error(),
		})
		return
	}
	result := s.Login(ctx)
	ctx.JSON(http.StatusOK, result)
}

func UserRegisterHandler(ctx *gin.Context) {
	var s service.UserRegisterService
	if err := ctx.ShouldBindJSON(&s); err != nil {
		ctx.JSON(http.StatusOK, serializer.Response{
			Code:  40001,
			Msg:   "操作失败，请检查参数",
			Error: err.Error(),
		})
		return
	}
	result := s.Register()
	ctx.JSON(http.StatusOK, result)
}

func UserLogoutHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.JSON(http.StatusOK, serializer.Response{
		Code: 0,
		Msg:  "退出登录",
	})
}

func UserUpdateHandler(ctx *gin.Context) {
	var service service.UserUpdateService
	if err := ctx.ShouldBindJSON(&service); err != nil {
		ctx.JSON(http.StatusOK, serializer.Response{
			Code:  40001,
			Msg:   "操作失败，请检查参数",
			Error: err.Error(),
		})
		return
	}
	user := middleware.GetLoginUser(ctx)
	service.UserId = user.Id
	result := service.Update()
	ctx.JSON(http.StatusOK, result)
}

func UserListHandler(ctx *gin.Context) {
	s := service.UserListService{}
	if !lib.RequestQueryParamBind(ctx, &s) {
		return
	}
	if s.PageSize == 0 {
		s.PageSize = 10
	}

	result := s.Query()
	ctx.JSON(http.StatusOK, result)
}

func AddUserHandler(ctx *gin.Context) {
	s := service.UserAddService{}
	if !lib.RequestParamBind(ctx, &s) {
		return
	}
	result := s.Add()
	ctx.JSON(http.StatusOK, result)
}

func ForbiddenHandler(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, nil)
}
