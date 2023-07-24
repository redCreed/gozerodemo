package logic

import (
	"context"
	"gozerodemo/app/reply/internal/svc"
	"gozerodemo/app/reply/reply"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentsLogic {
	return &CommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentsLogic) Comments(in *reply.CommentsRequest) (*reply.CommentsResponse, error) {

	return &reply.CommentsResponse{
		IsEnd:       true,
		CreatedTime: time.Now().Unix(),
	}, nil
}
