package project

type Header struct {
	ContentType string `header:"content-type" binding:"required"`
}

type AddReqParam struct {
	Group string `form:"group" binding:"required"`
}

type queryReqParam struct {
	GroupId string `form:"groupId""`
	PageNum int    `form:"pageNum"`
}
