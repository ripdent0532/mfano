package group

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rip0532/mfano/model"
)

func Register(group *gin.RouterGroup) {
	group.GET("/groups", groupList)
}

func groupList(context *gin.Context) {
	result := model.NewGroup().List()
	context.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
