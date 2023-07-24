package logic

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gozerodemo/app/product/internal/svc"
	"gozerodemo/app/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckAndUpdateStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckAndUpdateStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckAndUpdateStockLogic {
	return &CheckAndUpdateStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 秒杀从redis使用lua脚本进行判断是否有库存
const (
	luaCheckAndUpdateScript = `
local counts = redis.call("HMGET", KEYS[1], "total", "seckill")
local total = tonumber(counts[1])
local seckill = tonumber(counts[2])
if seckill + 1 <= total then
	redis.call("HINCRBY", KEYS[1], "seckill", 1)
	return 1
end
return 0
`
)

// CheckAndUpdateStock 检查库存
func (l *CheckAndUpdateStockLogic) CheckAndUpdateStock(in *product.CheckAndUpdateStockRequest) (*product.CheckAndUpdateStockResponse, error) {
	val, err := l.svcCtx.BizRedis.EvalCtx(l.ctx, luaCheckAndUpdateScript, []string{stockKey(in.ProductId)})
	if err != nil {
		return nil, err
	}
	if val.(int64) == 0 {
		return nil, status.Errorf(codes.ResourceExhausted, fmt.Sprintf("insufficient stock: %d", in.ProductId))
	}
	return &product.CheckAndUpdateStockResponse{}, nil

	return &product.CheckAndUpdateStockResponse{}, nil
}

func stockKey(pid int64) string {
	return fmt.Sprintf("stock:%d", pid)
}
