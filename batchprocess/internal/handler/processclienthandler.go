package handler

import (
	"net/http"

	"batchprocess/internal/logic"
	"batchprocess/internal/svc"
	"batchprocess/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ProcessClientHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProcessClientReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewProcessClientLogic(r.Context(), svcCtx)
		resp, err := l.ProcessClient(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
