package logic

import (
	"context"

	"gozerodemo/app/user/internal/svc"
	"gozerodemo/app/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserReceiveAddressListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserReceiveAddressListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserReceiveAddressListLogic {
	return &GetUserReceiveAddressListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取收获地址列表
func (l *GetUserReceiveAddressListLogic) GetUserReceiveAddressList(in *user.UserReceiveAddressListReq) (*user.UserReceiveAddressListRes, error) {
	// todo: add your logic here and delete this line

	return &user.UserReceiveAddressListRes{}, nil
}
