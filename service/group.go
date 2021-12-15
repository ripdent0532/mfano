package service

import (
	"time"

	"github.com/rip0532/mfano/model"
	"github.com/rip0532/mfano/serializer"
)

type GroupQueryService struct {
}

type groupQueryMapper struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Created   time.Time `json:"created"`
	Creator   string    `json:"creator"`
	CreatorId int       `json:"creator_id" db:"creator_id"`
}

func (service *GroupQueryService) Query() serializer.Response {
	var result []groupQueryMapper
	sql := `select * from mfano_group order by id asc`
	if err := model.DB.Select(&result, sql); err != nil {
		return serializer.Response{
			Code:  40001,
			Msg:   "查询分组失败",
			Error: err.Error(),
		}
	}
	return serializer.Response{
		Code: 0,
		Data: result,
	}

}
