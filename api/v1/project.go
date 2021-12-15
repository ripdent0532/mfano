package v1

import (
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rip0532/mfano/middleware"

	"github.com/rip0532/mfano/service"

	"github.com/gin-gonic/gin"
	"github.com/rip0532/mfano/lib"
	"github.com/rip0532/mfano/lib/archive"
	"github.com/rip0532/mfano/lib/constant"
	logger "github.com/rip0532/mfano/lib/log"
	"github.com/rip0532/mfano/serializer"
)

func ProjectQueryHandler(ctx *gin.Context) {
	service := service.ProjectQueryService{}
	if !lib.RequestQueryParamBind(ctx, &service) {
		return
	}
	result := service.Query()
	ctx.JSON(http.StatusOK, result)
}

func ProjectAddHandler(ctx *gin.Context) {
	files, err := upload(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, serializer.Response{
			Code:  40001,
			Msg:   "文件上传失败",
			Error: err.Error(),
		})
		return
	}

	zip := archive.NewZip()
	for _, file := range files {
		err := zip.UnZip(constant.DstDir, file)
		if err != nil {
			logger.Error.Printf("upload file error. %s\n", err)
			ctx.JSON(http.StatusOK, serializer.Response{
				Code:  40001,
				Msg:   "文件解压失败",
				Error: err.Error(),
			})
			// 清理文件
			lib.RemoveFileOrFolder(constant.DstDir + "/" + file.Timestamp)
			return
		}
	}
	service := service.ProjectAddService{}
	if !lib.RequestQueryParamBind(ctx, &service) {
		return
	}
	service.Files = files
	service.User = middleware.GetLoginUser(ctx)
	result := service.Add()
	ctx.JSON(http.StatusOK, result)
}

// 原始文件使用时间戳建立目录进行区隔
func upload(context *gin.Context) (files []archive.UploadFile, err error) {
	mediaType, params, _ := mime.ParseMediaType(context.Request.Header.Get("Content-Type"))
	if strings.HasPrefix(mediaType, "multipart/") {
		files = make([]archive.UploadFile, 0, 10)
		multiReader := multipart.NewReader(context.Request.Body, params["boundary"])
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		for {
			p, err := multiReader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			targetPath := constant.HomeDir + "/" + timestamp
			os.MkdirAll(constant.HomeDir+"/"+timestamp, os.ModePerm)
			filepath := targetPath + "/" + p.FileName()
			file, _ := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
			_, err = io.Copy(file, p)
			if err != nil {
				return nil, err
			}
			file.Close()
			files = append(files, archive.UploadFile{File: p.FileName(), Timestamp: timestamp})
		}
	}
	return files, nil
}
