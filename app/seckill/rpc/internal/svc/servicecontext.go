package svc

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gozerodemo/app/product/productclient"
	"gozerodemo/app/seckill/rpc/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	BizRedis    *redis.Redis
	ProductRPC  productclient.Product
	KafkaPusher *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
