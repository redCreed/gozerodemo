package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/mr"
	"gozerodemo/app/api/internal/svc"
	"gozerodemo/app/api/internal/types"
	"gozerodemo/app/product/product"
	"gozerodemo/app/reply/reply"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductDetailLogic {
	return &ProductDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductDetailLogic) ProductDetail(req *types.ProductDetailRequest) (resp *types.ProductDetailResponse, err error) {
	var (
		p  *product.ProductItem
		cs *reply.CommentsResponse
	)

	if err := mr.Finish(func() error {
		var err error
		t := time.Now()
		if p, err = l.svcCtx.ProductRPC.Product(l.ctx, &product.ProductItemRequest{ProductId: req.ProductID}); err != nil {
			return err
		}
		fmt.Println("ttt:", time.Since(t))
		return nil
	}, func() error {
		var err error
		if cs, err = l.svcCtx.ReplyRPC.Comments(l.ctx, &reply.CommentsRequest{TargetId: req.ProductID}); err != nil {
			logx.Errorf("get comments error: %v", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	var comments []*types.Comment
	for _, c := range cs.Comments {
		comments = append(comments, &types.Comment{
			ID:      c.Id,
			Content: c.Content,
		})
	}
	return &types.ProductDetailResponse{
		Product: &types.Product{
			ID:   p.ProductId,
			Name: p.Name,
		},
		Comments: comments,
	}, nil

}
