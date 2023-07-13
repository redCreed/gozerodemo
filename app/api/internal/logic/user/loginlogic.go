package user

import (
	"context"
	"fmt"
	"gozerodemo/app/user/user"

	"gozerodemo/app/api/internal/svc"
	"gozerodemo/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	//test处理
	in := &user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	ret, err := l.svcCtx.UserRPC.Login(context.Background(), in)
	fmt.Println("rpc resp:", ret.Id, ret.Phone, ret.Username)
	tt := &types.LoginResp{
		AccessToken:  "token",
		AccessExpire: 1,
	}
	return tt, err
}
