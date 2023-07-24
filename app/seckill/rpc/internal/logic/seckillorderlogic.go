package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/limit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gozerodemo/app/product/product"
	"gozerodemo/pkg/batcher"
	"strconv"
	"time"

	"gozerodemo/app/seckill/rpc/internal/svc"
	"gozerodemo/app/seckill/rpc/seckill"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	limitPeriod       = 10
	limitQuota        = 1
	seckillUserPrefix = "seckill#u#"
	localCacheExpire  = time.Second * 60

	batcherSize     = 100
	batcherBuffer   = 100
	batcherWorker   = 10
	batcherInterval = time.Second
)

type SeckillOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	limiter    *limit.PeriodLimit
	localCache *collection.Cache
	batcher    *batcher.Batcher
}

func NewSeckillOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillOrderLogic {
	localCache, err := collection.NewCache(localCacheExpire)
	if err != nil {
		panic(err)
	}
	s := &SeckillOrderLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		localCache: localCache,
		limiter:    limit.NewPeriodLimit(limitPeriod, limitQuota, svcCtx.BizRedis, seckillUserPrefix),
	}
	b := batcher.New(
		batcher.WithSize(batcherSize),
		batcher.WithBuffer(batcherBuffer),
		batcher.WithWorker(batcherWorker),
		batcher.WithInterval(batcherInterval),
	)

	b.Sharding = func(key string) int {
		pid, _ := strconv.ParseInt(key, 10, 64)
		return int(pid) % batcherWorker
	}
	b.Do = func(ctx context.Context, val map[string][]interface{}) {
		var msgs []*KafkaData
		for _, vs := range val {
			for _, v := range vs {
				msgs = append(msgs, v.(*KafkaData))
			}
		}
		kd, err := json.Marshal(msgs)
		if err != nil {
			logx.Errorf("Batcher.Do json.Marshal msgs: %v error: %v", msgs, err)
		}
		if err = s.svcCtx.KafkaPusher.Push(string(kd)); err != nil {
			logx.Errorf("KafkaPusher.Push kd: %s error: %v", string(kd), err)
		}
	}
	s.batcher = b
	s.batcher.Start()
	return s
}

func (l *SeckillOrderLogic) SeckillOrder(in *seckill.SeckillOrderRequest) (*seckill.SeckillOrderResponse, error) {
	code, _ := l.limiter.Take(strconv.FormatInt(in.UserId, 10))
	if code == limit.OverQuota {
		return nil, status.Errorf(codes.OutOfRange, "Number of requests exceeded the limit")
	}
	p, err := l.svcCtx.ProductRPC.Product(l.ctx, &product.ProductItemRequest{ProductId: in.ProductId})
	if err != nil {
		return nil, err
	}
	//校验库存
	if p.Stock <= 0 {
		return nil, status.Errorf(codes.OutOfRange, "Insufficient stock")
	}
	//抢购秒杀存到kafka中，由下游kafka进行消费进行dtm进行确定订单与减少库存
	if err = l.batcher.Add(strconv.FormatInt(in.ProductId, 10), &KafkaData{Uid: in.UserId, Pid: in.ProductId}); err != nil {
		logx.Errorf("l.batcher.Add uid: %d pid: %d error: %v", in.UserId, in.ProductId, err)
	}
	return &seckill.SeckillOrderResponse{}, nil
}
