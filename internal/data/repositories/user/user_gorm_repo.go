package user

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/newpurr/stock/internal/biz"
	"github.com/newpurr/stock/internal/data"
	"github.com/newpurr/stock/pkg/utils/gormutils"
	"github.com/newpurr/stock/pkg/utils/pagination"
	"gorm.io/gorm"
)

type userGormRepo struct {
	data.DefaultDataSource
	log *log.Helper
}

// NewUserGormRepo .
func NewUserGormRepo(data data.DefaultDataSource, logger log.Logger) biz.UserRepo {
	return &userGormRepo{
		DefaultDataSource: data,
		log:               log.NewHelper(logger),
	}
}

func (r *userGormRepo) CreateUser(ctx context.Context, c biz.UserDO) (result biz.UserDO, err error) {
	ent := NewUserPOFromDO(c)

	err = r.Gorm(ctx).Create(&ent).Error

	return ent.ToDO(), err
}

func (r *userGormRepo) UpdateUser(ctx context.Context, pk uint32, c biz.UserDO) (result biz.UserDO, err error) {
	if pk <= 0 {
		return biz.UserDO{}, errors.New("primary key must be greater than 0")
	}
	err = r.builder(ctx).whereId(pk).Save(&c).Error

	return c, err
}

func (r *userGormRepo) DeleteUser(ctx context.Context, pk uint32) (err error) {
	//return r.builder(ctx).whereId(pk).build().Update("is_deleted", 1).Error
	return r.Gorm(ctx).Delete(&UserPO{}, pk).Error
}

func (r *userGormRepo) GetUser(ctx context.Context, id uint32) (result biz.UserDO, err error) {
	err = r.Gorm(ctx).First(&result, id).Error

	return
}

func (r userGormRepo) buildBatchQueryBuilder(ctx context.Context, search biz.UserSearch) *gorm.DB {
	var (
		builder = r.builder(ctx)
	)
	if len(search.Ids) > 0 {
		// 主键筛选时，允许丢弃其他条件
		builder.resetWhere().whereIdIn(search.Ids)
		return builder.build()
	}

	return builder.build()
}

func (r *userGormRepo) BatchGetUser(ctx context.Context, vo biz.UserSearch) (result []biz.UserDO, err error) {
	var (
		gormTx = r.buildBatchQueryBuilder(ctx, vo)
	)

	err = gormTx.Scan(&result).Error

	return
}

func (r *userGormRepo) PagingQueryUser(ctx context.Context, vo biz.UserPagingSearch) (result []biz.UserDO, paginator pagination.Paginator, err error) {
	var (
		gormTx = r.buildBatchQueryBuilder(ctx, *vo.UserSearch)
	)

	err = gormutils.Paginate(
		gormTx,
		pagination.NewSearchPaginator(int64(vo.CurrentPage), int64(vo.PageSize)),
		func(db *gorm.DB, p pagination.Paginator) error {
			paginator = p
			return db.Scan(&result).Error
		},
	)

	return
}

func (r *userGormRepo) builder(ctx context.Context) *userQueryBuilder {
	return query(r.Gorm(ctx)).from(UserTableName)
}
