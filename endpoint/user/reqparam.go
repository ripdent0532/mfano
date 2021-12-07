package user

type ChangeUserInfoReqParam struct {
	Name            string `form:"username"`
	NickName        string `form:"nickname"`
	Email           string `form:"email"`
	OriginPassword  string `form:"originPassword"`
	NewPassword     string `form:"newPassword"`
	ConfirmPassword string `form:"confirmPassword"`
}

type AddUserReqParam struct {
	Name     string `form:"username" binding:"required"`
	NickName string `form:"nickname" binding:"required"`
	Email    string `form:"email" binding:"required"`
}

type QueryReqParam struct {
	UserName string `form:"userName"`
	PageNum  int    `form:"pageNum"`
}
