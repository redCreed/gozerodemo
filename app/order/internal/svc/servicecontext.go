package svc

import (
	"gozerodemo/app/order/internal/config"
	"gozerodemo/app/order/internal/model"
	"gozerodemo/app/product/productclient"
	"gozerodemo/app/user/userclient"
)

type ServiceContext struct {
	Config         config.Config
	OrderModel     model.OrdersModel
	OrderitemModel model.OrderitemModel
	ShippingModel  model.ShippingModel
	UserRpc        userclient.User
	ProductRpc     productclient.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
