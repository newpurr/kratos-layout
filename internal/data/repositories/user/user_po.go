package user

import (
	"github.com/newpurr/stock/internal/biz"
	"github.com/newpurr/stock/pkg/definition"
)

const (
	UserTableName = "kratos.user"
)

type (
	// UserPO 用户表
	UserPO struct {
		Id        uint32               `json:"id" gorm:"column:id"`                 // 主键
		Name      string               `json:"name" gorm:"column:name"`             // 用户名
		Pwd       string               `json:"pwd" gorm:"column:pwd"`               // 密码
		IsDeleted int32                `json:"is_deleted" gorm:"column:is_deleted"` // 是否删除 1:是  -1:否
		CreatedAt definition.LocalTime `json:"created_at" gorm:"column:created_at"` // 创建时间
		UpdatedAt definition.LocalTime `json:"updated_at" gorm:"column:updated_at"` // 更新时间
	}
)

// TableName UserPO 表名
func (UserPO) TableName() string {
	return UserTableName
}

func (po UserPO) ToDO() biz.UserDO {
	return biz.UserDO{
		Id:        po.Id,
		Name:      po.Name,
		Pwd:       po.Pwd,
		IsDeleted: po.IsDeleted,
		CreatedAt: po.CreatedAt.Time,
		UpdatedAt: po.UpdatedAt.Time,
	}
}

func NewUserPOFromDO(do biz.UserDO) UserPO {
	return UserPO{
		Id:        do.Id,
		Name:      do.Name,
		Pwd:       do.Pwd,
		IsDeleted: do.IsDeleted,
		CreatedAt: definition.LocalTime{Time: do.CreatedAt},
		UpdatedAt: definition.LocalTime{Time: do.UpdatedAt},
	}
}
