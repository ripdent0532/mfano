package login

// Form 登录参数
type Form struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
