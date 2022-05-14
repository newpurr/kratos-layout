package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/newpurr/stock/internal/biz"
	"github.com/newpurr/stock/internal/data"
	"time"
)

type userCacheRepo struct {
	data.DefaultDataSource
	biz.UserRepo
	Logger *log.Helper
}

// NewUserCacheRepo .
// 其实这里只依赖redis即可?????
func NewUserCacheRepo(data data.DefaultDataSource, logger log.Logger) biz.UserRepo {
	return &userCacheRepo{
		data,
		NewUserGormRepo(data, logger),
		log.NewHelper(logger),
	}
}

func (r *userCacheRepo) UpdateUser(ctx context.Context, pk uint32, g biz.UserDO) (channel biz.UserDO, err error) {
	channel, err = r.UserRepo.UpdateUser(ctx, pk, g)
	if err != nil {
		// 删除缓存，需要事务完成后去执行
		cacheKey := fmt.Sprintf("user:%d", pk)
		deleteCache := func(ctx context.Context, pk uint32, r *userCacheRepo) error {
			return r.RedisClient(ctx).Del(ctx, cacheKey).Err()
		}
		err = deleteCache(ctx, pk, r)
		if err != nil {
			return channel, err
		}

		//time.AfterFunc(time.Second*5, func() {
		//	err := deleteCache(context.Background(), pk, r)
		//	if err != nil {
		//		r.Logger.Error("delete %s fail.", cacheKey)
		//	}
		//})
	}
	return channel, err
}

func (r *userCacheRepo) GetUser(ctx context.Context, id uint32) (biz.UserDO, error) {
	var (
		ch       biz.UserDO // TODO 可以用作泛型T
		cacheKey = fmt.Sprintf("user:%d", id)
		redis    = r.RedisClient(ctx)
	)

	// 判断缓存是否存在，如果存在直接返回
	result, err := redis.Get(ctx, cacheKey).Result()
	if err == nil && len(result) > 0 {
		err := json.Unmarshal([]byte(result), &ch)
		if err != nil {
			goto QueryFromDb
		}

		// 序列化字符串到对象，然后返回
		return ch, err
	}

QueryFromDb:

	// 如果不存在直接从db获取，然后更新到数据库
	ch, err = r.UserRepo.GetUser(ctx, id)
	if err != nil {
		return ch, err
	}

	// TODO 事务完成后去执行
	bytes, _ := json.Marshal(ch)
	err = redis.Set(ctx, cacheKey, string(bytes), time.Hour).Err()
	if err != nil {
		// 记录日志，不更新缓存
		return ch, err
	}

	return ch, err
}
