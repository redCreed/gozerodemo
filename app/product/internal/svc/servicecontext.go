package svc

import (
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/sync/singleflight"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gozerodemo/app/product/internal/config"
	"gozerodemo/app/product/internal/model"
	"time"
)

type ServiceContext struct {
	Config   config.Config
	BizRedis *redis.Redis
	//可以使用gorm
	orm        *gorm.DB
	LocalCache *collection.Cache
	//防止缓存击穿
	SingleGroup singleflight.Group

	ProductModel   model.ProductModel
	CategoryModel  model.CategoryModel
	OperationModel model.ProductOperationModel
}

const localCacheExpire = time.Duration(time.Second * 60)

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	//本地缓存
	localCache, err := collection.NewCache(localCacheExpire)
	if err != nil {
		panic(err)
	}
	dsn := c.DataSource
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:         c,
		BizRedis:       redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		LocalCache:     localCache,
		orm:            db,
		ProductModel:   model.NewProductModel(conn, c.CacheRedis),
		CategoryModel:  model.NewCategoryModel(conn, c.CacheRedis),
		OperationModel: model.NewProductOperationModel(conn, c.CacheRedis),
	}
}
