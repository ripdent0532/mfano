package lib

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rip0532/mfano/serializer"
)

func RequestParamBind(ctx *gin.Context, dist interface{}) bool {
	if err := ctx.ShouldBindJSON(dist); err != nil {
		ctx.JSON(http.StatusOK, serializer.Response{
			Code:  40001,
			Msg:   "操作失败，请检查参数",
			Error: err.Error(),
		})
		return false
	}
	return true
}

func RequestQueryParamBind(ctx *gin.Context, dist interface{}) bool {
	if err := ctx.ShouldBindQuery(dist); err != nil {
		ctx.JSON(http.StatusOK, serializer.Response{
			Code:  40001,
			Msg:   "操作失败，请检查参数",
			Error: err.Error(),
		})
		return false
	}
	handlerPage(dist)
	return true
}

func handlerPage(object interface{}) {
	objType := reflect.TypeOf(object)
	pageSizeField, present := objType.Elem().FieldByName("PageSize")
	if present {
		structTag := pageSizeField.Tag
		if pageTag := structTag.Get("page"); pageTag != "" {
			pageSize := strings.Split(pageTag, "=")[1]
			handlerPageValue(object, pageSize)
		}

	}
}

func handlerPageValue(object interface{}, value string) {
	objValue := reflect.ValueOf(object)
	fieldValue := objValue.Elem().FieldByName("PageSize")
	if fieldValue.Int() == 0 {
		v, _ := strconv.Atoi(value)
		fieldValue.Set(reflect.ValueOf(v))
	}
}

func Md5File(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		return ""
	}
	defer file.Close()
	m := md5.New()
	_, err = io.Copy(m, file)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(m.Sum(nil))
}
