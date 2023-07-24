package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gozerodemo/app/api/internal/svc"
	"gozerodemo/app/api/internal/types"
	"gozerodemo/app/product/product"
	"time"
)

type RecommendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendLogic {
	return &RecommendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendLogic) Recommend(req *types.RecommendRequest) (resp *types.RecommendResponse, err error) {
	ll := &product.ProductListRequest{
		Cursor: req.Cursor,
		Ps:     int32(req.Ps),
	}
	res, err := l.svcCtx.ProductRPC.ProductList(context.Background(), ll)
	if err != nil {
		return nil, err
	}

	resp = new(types.RecommendResponse)
	resp.RecommendTime = time.Now().Unix()
	resp.IsEnd = false
	product := make([]*types.Product, 0)
	for _, v := range res.Products {
		product = append(product, &types.Product{
			ID:          v.ProductId,
			Name:        v.Name,
			Images:      nil,
			Description: v.Description,
			Price:       v.Price,
			Stock:       v.Stock,
			Category:    "",
			Status:      v.Status,
			CreateTime:  v.CreateTime,
			UpdateTime:  0,
		})
	}
	resp.Products = product
	return
}
