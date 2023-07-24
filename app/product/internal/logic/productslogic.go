package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gozerodemo/app/product/internal/svc"
	"gozerodemo/app/product/product"
)

type ProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductsLogic {
	return &ProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type T interface {
	string
}

type U interface{}

type V interface{}

func (l *ProductsLogic) Products(in *product.ProductRequest) (*product.ProductResponse, error) {
	//products := make(map[int64]*product.ProductItem)
	//pids := strings.Split(in.ProductIds, ",")
	////todo
	//ps, err := mr.MapReduce(func(source chan<- interface{}) {
	//	for _, pid := range pids {
	//		source <- pid
	//	}
	//}, func(item interface{}, writer mr.Writer[U], cancel func(error)) {
	//	pidStr := item.(string)
	//	pid, err := strconv.ParseInt(pidStr, 10, 64)
	//	if err != nil {
	//		return
	//	}
	//	p, err := l.svcCtx.ProductModel.FindOne(l.ctx, pid)
	//	if err != nil {
	//		return
	//	}
	//	writer.Write(p)
	//}, func(pipe <-chan U, writer mr.Writer[V], cancel func(error)) {
	//	var r []*model.Product
	//	for p := range pipe {
	//		r = append(r, p.(*model.Product))
	//	}
	//	writer.Write(r)
	//})
	//if err != nil {
	//	return nil, err
	//}
	//for _, p := range ps.([]*model.Product) {
	//	products[p.Id] = &product.ProductItem{
	//		ProductId: p.Id,
	//		Name:      p.Name,
	//	}
	//}
	//return &product.ProductResponse{Products: products}, nil
	return nil, nil

}
