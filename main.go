package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	logger "github.com/rip0532/mfano/lib/log"

	"github.com/rip0532/mfano/lib/db"

	"github.com/rip0532/mfano/endpoint/user"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/rip0532/mfano/endpoint/group"
	"github.com/rip0532/mfano/endpoint/project"

	"github.com/gin-gonic/gin"
	translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/rip0532/mfano/endpoint/login"
	"github.com/rip0532/mfano/lib"
	"github.com/rip0532/mfano/lib/constant"
	"github.com/rip0532/mfano/middleware"
	"golang.org/x/sync/errgroup"
)

var (
	g     errgroup.Group
	trans ut.Translator
)

func init() {
	if !lib.FolderExists(constant.DstDir) {
		os.MkdirAll(constant.DstDir, os.ModeDir)
		logger.Info.Printf("创建文件夹：%s\n", constant.DstDir)
	}
	if !lib.FolderExists(constant.Db_Dir) {
		os.MkdirAll(constant.Db_Dir, os.ModeDir)
		logger.Info.Printf("创建文件夹：%s\n", constant.Db_Dir)
	}
	// 连接数据库
	db.Open()
}

func main() {
	if err := InitializeTrans(); err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	gin.SetMode(gin.ReleaseMode)
	resourceServer := &http.Server{
		Addr:         ":8081",
		Handler:      staticResourcesRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	viewServer := &http.Server{
		Addr:         ":8082",
		Handler:      viewRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	apiServer := &http.Server{
		Addr:         ":8080",
		Handler:      serverRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return resourceServer.ListenAndServe()
	})

	g.Go(func() error {
		return viewServer.ListenAndServe()
	})

	g.Go(func() error {
		return apiServer.ListenAndServe()
	})

	logger.Info.Println("Mfano已启动")

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func staticResourcesRouter() http.Handler {
	logger.Info.Println("初始化静态资源服务")
	e := gin.Default()
	e.Use(middleware.Session())
	e.NoRoute(func(context *gin.Context) {
		context.Writer.WriteString("<head><title>404 - Mfano</title></head>" +
			"<body><div style='margin: 0 auto; width: 60%; text-align: center;'><h2>糟糕！页面不见了！！！🛸</h2></div></body>")
	})
	authorized := e.Group("/", middleware.SessionHandler())
	authorized.StaticFS("/", http.Dir(constant.DstDir))
	return e
}

func viewRouter() http.Handler {
	logger.Info.Println("初始化网页资源服务")
	e := gin.Default()
	e.Static("/", "./views")
	return e
}

func serverRouter() http.Handler {
	logger.Info.Println("初始化API接口服务")
	e := gin.Default()
	e.Use(middleware.Session(), middleware.Cros())
	authorized := e.Group("/", middleware.SessionHandler())
	login.Register(e, authorized)
	project.Register(authorized)
	user.Register(authorized)
	group.Register(authorized)
	return e
}

func InitializeTrans() (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := field.Tag.Get("json")
			return name
		})
		zhT := zh.New()
		uni := ut.New(zhT, zhT)
		trans, _ = uni.GetTranslator("zh")
		err = translations.RegisterDefaultTranslations(v, trans)
		return
	}
	return
}
