package logic

import (
	"context"
	"gozerodemo/app/seckill/rpc/internal/svc"
	"gozerodemo/app/seckill/rpc/seckill"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	logx.Logger
}
type KafkaData struct {
	Uid int64 `json:"uid"`
	Pid int64 `json:"pid"`
}

func NewSeckillProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillProductsLogic {
	return &SeckillProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SeckillProductsLogic) SeckillProducts(in *seckill.SeckillProductsRequest) (*seckill.SeckillProductsResponse, error) {

	return &seckill.SeckillProductsResponse{}, nil
}
