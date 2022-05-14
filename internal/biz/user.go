package biz

import (
	"context"
	"github.com/asaskevich/EventBus"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/newpurr/stock/pkg/utils/pagination"
	"time"
)

type (
	// UserDO 用户表
	UserDO struct {
		Id        uint32    `json:"id"`         // 主键
		Name      string    `json:"name"`       // 用户名
		Pwd       string    `json:"pwd"`        // 密码
		IsDeleted int32     `json:"is_deleted"` // 是否删除 1:是  -1:否
		CreatedAt time.Time `json:"created_at"` // 创建时间
		UpdatedAt time.Time `json:"updated_at"` // 更新时间
	}

	// UserSearch 用户表筛选
	UserSearch struct {
		Ids []uint32 `json:"id"` // 主键
	}

	// UserPagingSearch 用户表分页筛选
	UserPagingSearch struct {
		*UserSearch
		*PaginateReq
	}
)

type UserRepo interface {
	CreateUser(context.Context, UserDO) (UserDO, error)
	UpdateUser(context.Context, uint32, UserDO) (UserDO, error)
	GetUser(context.Context, uint32) (UserDO, error)
	DeleteUser(context.Context, uint32) (err error)
	BatchGetUser(context.Context, UserSearch) (result []UserDO, err error)
	PagingQueryUser(context.Context, UserPagingSearch) (result []UserDO, paginator pagination.Paginator, err error)
}

type UserUsecase struct {
	tm   TransactionManager
	bus  EventBus.BusPublisher
	log  *log.Helper
	repo UserRepo
}

func NewUserUsecase(tm TransactionManager, bus EventBus.BusPublisher, logger log.Logger, repo UserRepo) *UserUsecase {
	return &UserUsecase{tm: tm, log: log.NewHelper(logger), bus: bus, repo: repo}
}

func (uc *UserUsecase) Create(ctx context.Context, g UserDO) (ent UserDO, err error) {
	defer func() {
		if err == nil {
			uc.bus.Publish("ent:user:create", ent)
		}
	}()
	_ = uc.tm.Transaction(ctx, func(ctx context.Context) error {
		ent, err = uc.repo.CreateUser(ctx, g)
		return err
	})
	return
}

func (uc *UserUsecase) Update(ctx context.Context, g UserDO) (ent UserDO, err error) {
	defer func() {
		if err == nil {
			uc.bus.Publish("ent:user:update", ent)
		}
	}()
	_ = uc.tm.Transaction(ctx, func(ctx context.Context) error {
		ent, err = uc.repo.UpdateUser(ctx, g.Id, g)
		return err
	})
	return
}

func (uc *UserUsecase) Get(ctx context.Context, pk uint32) (UserDO, error) {
	return uc.repo.GetUser(ctx, pk)
}

func (uc *UserUsecase) PagingQueryUser(ctx context.Context, vo UserPagingSearch) (result []UserDO, paginator pagination.Paginator, err error) {
	return uc.repo.PagingQueryUser(ctx, vo)
}

func (uc *UserUsecase) Delete(ctx context.Context, pk uint32) (err error) {
	defer func() {
		if err == nil {
			uc.bus.Publish("ent:user:delete", pk)
		}
	}()
	_ = uc.tm.Transaction(ctx, func(ctx context.Context) error {
		err = uc.repo.DeleteUser(ctx, pk)
		return err
	})
	return
}
