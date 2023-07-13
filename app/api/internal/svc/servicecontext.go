package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gozerodemo/app/api/internal/config"
	"gozerodemo/app/user/userclient"
)

type ServiceContext struct {
	Config  config.Config
	UserRPC userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRPC: userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
	}
}
