package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/threading"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gozerodemo/app/product/internal/model"
	"gozerodemo/app/product/internal/svc"
	"gozerodemo/app/product/product"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	defaultPageSize = 10
	defaultLimit    = 300
	expireTime      = 3600 * 24 * 3
)

type ProductListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductListLogic {
	return &ProductListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductListLogic) ProductList(in *product.ProductListRequest) (*product.ProductListResponse, error) {
	in.CategoryId = 8
	//先从cache查询，查不到就从db查询  model层缓存默认7天
	_, err := l.svcCtx.CategoryModel.FindOne(l.ctx, int64(in.CategoryId))
	if err == model.ErrNotFound {

		return nil, status.Error(codes.NotFound, "category not found")
	}
	if in.Cursor == 0 {
		in.Cursor = time.Now().Unix()
	}
	if in.Ps == 0 {
		in.Ps = defaultPageSize
	}
	var (
		isCache, isEnd   bool
		lastID, lastTime int64
		firstPage        []*product.ProductItem
		products         []*model.Product
	)
	//从缓存zset中查询  这里调用查询缓存方法忽略了error，为什么要忽略这个error呢，
	//因为我们期望的是尽最大可能的给用户返回数据，也就是redis挂掉了的话那我们就会从数据库查询数据返回给用户，
	//而不会因为redis挂掉而返回错误。
	pids, _ := l.cacheProductList(l.ctx, in.CategoryId, in.Cursor, int64(in.Ps))
	//orted Set中设置一个结束标志符，该标志符的member为-1，score为0
	if len(pids) == int(in.Ps) {
		isCache = true
		//缓存中查出数据后，需要判断数据的最后一条是否为-1，如果为-1的话说明列表已经加载到最后一页
		if pids[len(pids)-1] == -1 {
			isEnd = true
		}
		//使用mapreducer来并发查询
		products, err := l.productsByIds(l.ctx, pids)
		if err != nil {
			return nil, err
		}
		for _, p := range products {
			firstPage = append(firstPage, &product.ProductItem{
				ProductId:  p.Id,
				Name:       p.Name,
				CreateTime: p.CreateTime.Unix(),
			})
		}
	} else {
		var (
			err   error
			ctime = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")
		)
		products, err = l.svcCtx.ProductModel.CategoryProducts(l.ctx, ctime, int64(in.CategoryId), defaultLimit)
		if err != nil {
			return nil, err
		}
		fmt.Println("dddd:", len(products))
		var firstPageProducts []*model.Product
		if len(products) > int(in.Ps) {
			firstPageProducts = products[:int(in.Ps)]
		} else {
			firstPageProducts = products
			isEnd = true
		}
		for _, p := range firstPageProducts {
			firstPage = append(firstPage, &product.ProductItem{
				ProductId:  p.Id,
				Name:       p.Name,
				CreateTime: p.CreateTime.Unix(),
			})
		}
	}
	if len(firstPage) > 0 {
		pageLast := firstPage[len(firstPage)-1]
		lastID = pageLast.ProductId
		lastTime = pageLast.CreateTime
		if lastTime < 0 {
			lastTime = 0
		}
		for k, p := range firstPage {
			if p.CreateTime == in.Cursor && p.ProductId == in.ProductId {
				firstPage = firstPage[k:]
				break
			}
		}
	}
	ret := &product.ProductListResponse{
		IsEnd:     isEnd,
		Timestamp: lastTime,
		ProductId: lastID,
		Products:  firstPage,
	}
	if !isCache {
		threading.GoSafe(func() {
			if len(products) < defaultLimit && len(products) > 0 {
				endTime, _ := time.Parse("2006-01-02 15:04:05", "0000-00-00 00:00:00")
				products = append(products, &model.Product{Id: -1, CreateTime: endTime})
			}
			_ = l.addCacheProductList(context.Background(), products)
		})
	}
	return ret, nil
}

func (l *ProductListLogic) cacheProductList(ctx context.Context, cid int32, cursor, ps int64) ([]int64, error) {
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, categoryKey(cid), cursor, 0, 0, int(ps))
	if err != nil {
		return nil, err
	}
	var ids []int64
	for _, pair := range pairs {
		id, _ := strconv.ParseInt(pair.Key, 10, 64)
		ids = append(ids, id)
	}
	return ids, nil
}

func (l *ProductListLogic) addCacheProductList(ctx context.Context, products []*model.Product) error {
	if len(products) == 0 {
		return nil
	}
	for _, p := range products {
		score := p.CreateTime.Unix()
		if score < 0 {
			score = 0
		}
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, categoryKey(int32(p.Cateid)), score, strconv.Itoa(int(p.Id)))
		if err != nil {
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, categoryKey(int32(products[0].Cateid)), expireTime)
}

func categoryKey(cid int32) string {
	return fmt.Sprintf("category:%d", cid)
}

func (l *ProductListLogic) productsByIds(ctx context.Context, pids []int64) ([]*model.Product, error) {
	//products, err := mr.MapReduce(func(source chan<- interface{}) {
	//	for _, pid := range pids {
	//		source <- pid
	//	}
	//}, func(item interface{}, writer mr.Writer, cancel func(error)) {
	//	pid := item.(int64)
	//	p, err := l.svcCtx.ProductModel.FindOne(ctx, pid)
	//	if err != nil {
	//		cancel(err)
	//		return
	//	}
	//	writer.Write(p)
	//}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
	//	var ps []*model.Product
	//	for item := range pipe {
	//		p := item.(*model.Product)
	//		ps = append(ps, p)
	//	}
	//	writer.Write(ps)
	//})
	//if err != nil {
	//	return nil, err
	//}
	//return products.([]*model.Product), nil
	return nil, nil
}
