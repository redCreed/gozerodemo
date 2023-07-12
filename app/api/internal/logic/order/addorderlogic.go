package order

import (
	"context"

	"gozerodemo/app/api/internal/svc"
	"gozerodemo/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddOrderLogic {
	return &AddOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddOrderLogic) AddOrder(req *types.OrderAddReq) (resp *types.OrderAddResp, err error) {
	// todo: add your logic here and delete this line

	return
}
