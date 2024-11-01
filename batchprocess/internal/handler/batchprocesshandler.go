package handler

import (
	"net/http"

	"batchprocess/internal/logic"
	"batchprocess/internal/svc"
	"batchprocess/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func BatchProcessHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchProcessReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewBatchProcessLogic(r.Context(), svcCtx)
		resp, err := l.BatchProcess(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
