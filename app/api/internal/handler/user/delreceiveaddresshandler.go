package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozerodemo/app/api/internal/logic/user"
	"gozerodemo/app/api/internal/svc"
	"gozerodemo/app/api/internal/types"
)

func DelReceiveAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserReceiveAddressDelReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewDelReceiveAddressLogic(r.Context(), svcCtx)
		resp, err := l.DelReceiveAddress(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
