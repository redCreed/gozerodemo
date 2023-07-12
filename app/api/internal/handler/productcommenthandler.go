package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozerodemo/app/api/internal/logic"
	"gozerodemo/app/api/internal/svc"
	"gozerodemo/app/api/internal/types"
)

func ProductCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProductCommentRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewProductCommentLogic(r.Context(), svcCtx)
		resp, err := l.ProductComment(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
