package logic

import (
	"context"
	"github.com/dtm-labs/dtmcli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gozerodemo/app/product/internal/svc"
	"gozerodemo/app/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckProductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckProductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckProductStockLogic {
	return &CheckProductStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckProductStockLogic) CheckProductStock(in *product.UpdateProductStockRequest) (*product.UpdateProductStockResponse, error) {
	p, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.ProductId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if p.Stock < in.Num {
		return nil, status.Error(codes.ResourceExhausted, dtmcli.ResultFailure)
	}

	return &product.UpdateProductStockResponse{}, nil
}
