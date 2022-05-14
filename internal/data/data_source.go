package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/newpurr/stock/internal/biz"
	"github.com/newpurr/stock/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	slog "log"
	"os"
	"time"
)

/*--------------------------------------------------------------------------
| 定义DataSource接口
|-------------------------------------------------------------------------- */
type (
	TransactionDataSource interface {
		biz.TransactionManager
	}
	GormDataSource interface {
		Gorm(ctx context.Context) *gorm.DB
	}
	RedisDataSource interface {
		RedisClient(ctx context.Context) *redis.Client
	}
	DefaultDataSource interface {
		TransactionDataSource
		RedisDataSource
		GormDataSource
	}
)

/*--------------------------------------------------------------------------
| 定义DataSource默认实现
|-------------------------------------------------------------------------- */
type (
	DefaultDataSourceImpl struct {
		config      *conf.Data
		gorm        *gorm.DB
		redisClient *redis.Client
	}
	transactionSession struct{}
)

func (d *DefaultDataSourceImpl) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sourceImpl := &DefaultDataSourceImpl{
			gorm:        tx,
			redisClient: d.redisClient.WithContext(ctx),
		}
		ctx = injectSessionOnce(ctx, sourceImpl)

		return fn(ctx)
	})

}

func (d *DefaultDataSourceImpl) Gorm(ctx context.Context) *gorm.DB {
	session, ok := loadSession(ctx).(*DefaultDataSourceImpl)
	if ok {
		return session.gorm
	}

	return d.gorm
}

func (d *DefaultDataSourceImpl) RedisClient(ctx context.Context) *redis.Client {
	_ = ctx
	return d.redisClient
}

/*--------------------------------------------------------------------------
| help function
|-------------------------------------------------------------------------- */
func injectSessionOnce(ctx context.Context, d *DefaultDataSourceImpl) context.Context {
	v := ctx.Value(transactionSession{})
	if v == nil {
		ctx = context.WithValue(ctx, transactionSession{}, d)
	}

	return ctx
}

func loadSession(ctx context.Context) DefaultDataSource {
	DataSource, ok := ctx.Value(transactionSession{}).(*DefaultDataSourceImpl)
	if !ok {
		return nil
	}
	return DataSource
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB, rdb *redis.Client) (DefaultDataSource, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	dataImpl := DefaultDataSourceImpl{
		config:      c,
		gorm:        db,
		redisClient: rdb,
	}
	return &dataImpl, cleanup, nil
}

func NewTransactionDataSource(source DefaultDataSource) biz.TransactionManager {
	return source
}

/*--------------------------------------------------------------------------
| gorm实现和redis client实现
|-------------------------------------------------------------------------- */

// NewDB .
func NewDB(c *conf.Data) *gorm.DB {
	// 终端打印输入 sql 执行记录
	newLogger := logger.New(
		slog.New(os.Stdout, "\r\n", slog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢查询 SQL 阈值
			Colorful:      true,        // 禁用彩色打印
			//IgnoreRecordNotFoundError: false,
			LogLevel: logger.Info, // Log lever
		},
	)

	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名是否加 s
		},
	})

	if err != nil {
		log.Errorf("failed opening connection to sqlite: %v", err)
		panic("failed to connect database")
	}

	return db
}

func NewRedis(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		PoolSize:     100,
	})
	//redisClient.AddHook(redisotel.TracingHook{})
	//if err := redisClient.Close(); err != nil {
	//	log.Error(err)
	//}
	return rdb
}
