package user

import (
	"context"

	"gozerodemo/app/api/internal/svc"
	"gozerodemo/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCollectionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCollectionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCollectionListLogic {
	return &UserCollectionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCollectionListLogic) UserCollectionList(req *types.UserCollectionListReq) (resp *types.UserCollectionListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
