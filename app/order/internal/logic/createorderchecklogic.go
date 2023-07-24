package logic

import (
	"context"
	"gozerodemo/app/order/internal/svc"
	"gozerodemo/app/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderCheckLogic {
	return &CreateOrderCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderCheckLogic) CreateOrderCheck(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {

	return &order.CreateOrderResponse{}, nil
}
