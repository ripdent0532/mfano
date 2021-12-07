# Mfano

> 起源于原型托管需求，用轻量便捷的方式管理和分享原型设计！

### 环境变量

| 变量             | 解释                                    |
| ---------------- | --------------------------------------- |
| HOME_DIR         | 原始文件保存位置                        |
| DST_DIR          | 解压文件保存位置                        |
| DB_DIR           | 数据库文件位置（sqlite3）               |
| STATIC_HOST      | 静态文件服务域名                        |
| API_HOST         | API接口服务域名                         |
| DOMAIN           | 域名。用于子域名间共享cookie            |
| DEFAULT_PASSWORD | 默认密码                                |
| MODE             | 模式。DEV（开发模式）、PROD（生产模式） |

### 构建
已提供编写好的`Dockerfile`