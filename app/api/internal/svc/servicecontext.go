package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gozerodemo/app/api/internal/config"
	"gozerodemo/app/product/productclient"
	"gozerodemo/app/reply/replyclient"
	"gozerodemo/app/user/userclient"
)

type ServiceContext struct {
	Config     config.Config
	UserRPC    userclient.User
	ProductRPC productclient.Product
	ReplyRPC   replyclient.Reply
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		UserRPC:    userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
		ProductRPC: productclient.NewProduct(zrpc.MustNewClient(c.ProductRPC)),
		ReplyRPC:   replyclient.NewReply(zrpc.MustNewClient(c.ReplyRPC)),
	}
}
