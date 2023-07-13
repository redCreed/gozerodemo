package logic

import (
	"context"

	"gozerodemo/app/user/internal/svc"
	"gozerodemo/app/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserReceiveAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserReceiveAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserReceiveAddressLogic {
	return &AddUserReceiveAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 添加收获地址
func (l *AddUserReceiveAddressLogic) AddUserReceiveAddress(in *user.UserReceiveAddressAddReq) (*user.UserReceiveAddressAddRes, error) {
	// todo: add your logic here and delete this line

	return &user.UserReceiveAddressAddRes{}, nil
}
