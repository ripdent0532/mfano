package constant

import "os"

// HomeDir 主目录
var HomeDir = os.Getenv("HOME_DIR")

// DstDir 文件解压目录
var DstDir = os.Getenv("DST_DIR")

// Db_Dir 数据库路径
var Db_Dir = os.Getenv("DB_DIR")
var StaticHost = os.Getenv("STATIC_HOST")
var DefaultPassword = os.Getenv("DEFAULT_PASSWORD")

// Domain 域名，多域名共享cookie，配置为：.domain.com
var Domain = os.Getenv("DOMAIN")

// Mode 运行模式，DEV
var Mode = os.Getenv("MODE")

var ApiHost = os.Getenv("API_HOST")
