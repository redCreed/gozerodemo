package logic

import (
	"context"

	"gozerodemo/app/order/internal/svc"
	"gozerodemo/app/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderDTMRevertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderDTMRevertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderDTMRevertLogic {
	return &CreateOrderDTMRevertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderDTMRevertLogic) CreateOrderDTMRevert(in *order.AddOrderReq) (*order.AddOrderResp, error) {
	// todo: add your logic here and delete this line

	return &order.AddOrderResp{}, nil
}
