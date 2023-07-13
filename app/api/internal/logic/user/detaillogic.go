package user

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"gozerodemo/app/user/user"

	"gozerodemo/app/api/internal/svc"
	"gozerodemo/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	userInfo, err := l.svcCtx.UserRPC.UserInfo(l.ctx, &user.UserInfoRequest{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}
	var user types.UserInfo
	//从一个结构体拷贝到另外一个结构体
	_ = copier.Copy(&user, userInfo.User)
	return &types.UserInfoResp{
		UserInfo: user,
	}, nil
}
