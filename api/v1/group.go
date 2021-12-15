package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rip0532/mfano/service"
)

func GroupQueryHandler(ctx *gin.Context) {
	s := service.GroupQueryService{}
	result := s.Query()
	ctx.JSON(http.StatusOK, result)
}
