package middleware

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/rip0532/mfano/lib/constant"
)

type SessionUser struct {
	Id        int64
	Name      string
	NickName  string
	AvatarURL string
	Email     string
}

const secret = "secret"
const sessionName = "session"

func Session() gin.HandlerFunc {
	gob.Register(SessionUser{})
	store := cookie.NewStore([]byte(secret))
	if constant.Mode != "DEV" {
		store.Options(sessions.Options{
			Path:     "/",
			Domain:   constant.Domain,
			MaxAge:   3600,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
	}
	return sessions.Sessions(sessionName, store)
}

func SessionHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, login := CheckSession(ctx)
		if !login {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func CheckSession(context *gin.Context) (interface{}, bool) {
	session := sessions.Default(context)
	user := session.Get("user")
	return user, !(user == nil)
}

func GetLoginUser(context *gin.Context) SessionUser {
	session := sessions.Default(context)
	user := session.Get("user")
	return user.(SessionUser)
}
