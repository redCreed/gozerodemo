package logic

import (
	"context"

	"gozerodemo/app/product/internal/svc"
	"gozerodemo/app/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperationProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	productListLogic *ProductListLogic
}

func NewOperationProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationProductsLogic {
	return &OperationProductsLogic{
		ctx:              ctx,
		svcCtx:           svcCtx,
		Logger:           logx.WithContext(ctx),
		productListLogic: NewProductListLogic(ctx, svcCtx),
	}
}

const (
	validStatus          = 1
	operationProductsKey = "operation#products"
)

func (l *OperationProductsLogic) OperationProducts(in *product.OperationProductsRequest) (*product.OperationProductsResponse, error) {
	//Cache来实现本地缓存的功能
	//先从本地缓存中查找，如果命中缓存则直接返回。没有命中缓存的话需要先从数据库中查询运营位商品id，
	//然后再聚合商品信息，最后回塞到本地缓存中
	//这种办法是用来解决热点问题
	opProducts, ok := l.svcCtx.LocalCache.Get(operationProductsKey)
	if ok {
		return &product.OperationProductsResponse{Products: opProducts.([]*product.ProductItem)}, nil
	}
	pos, err := l.svcCtx.OperationModel.OperationProducts(l.ctx, validStatus)
	if err != nil {
		return nil, err
	}
	var pids []int64
	for _, p := range pos {
		pids = append(pids, p.ProductId)
	}
	products, err := l.productListLogic.productsByIds(l.ctx, pids)
	if err != nil {
		return nil, err
	}
	var pItems []*product.ProductItem
	for _, p := range products {
		pItems = append(pItems, &product.ProductItem{
			ProductId: p.Id,
			Name:      p.Name,
		})
	}
	l.svcCtx.LocalCache.Set(operationProductsKey, pItems)
	return &product.OperationProductsResponse{Products: pItems}, nil

}
