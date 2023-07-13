package logic

import (
	"context"

	"gozerodemo/app/user/internal/svc"
	"gozerodemo/app/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditUserReceiveAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditUserReceiveAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserReceiveAddressLogic {
	return &EditUserReceiveAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 编辑收获地址
func (l *EditUserReceiveAddressLogic) EditUserReceiveAddress(in *user.UserReceiveAddressEditReq) (*user.UserReceiveAddressEditRes, error) {
	// todo: add your logic here and delete this line

	return &user.UserReceiveAddressEditRes{}, nil
}
