package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozerodemo/app/api/internal/logic"
	"gozerodemo/app/api/internal/svc"
	"gozerodemo/app/api/internal/types"
)

func ProductDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProductDetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewProductDetailLogic(r.Context(), svcCtx)
		resp, err := l.ProductDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
