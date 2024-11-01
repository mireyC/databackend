package geoContainer

import (
	"context"
	"fmt"
	"geoserver/api/internal/util"
	"net/http"

	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelDataSourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelDataSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelDataSourceLogic {
	return &DelDataSourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelDataSourceLogic) DelDataSource(req *types.DelDataSourceReq) (*types.ErrorResponse, error) {
	storeName := "bev" + "_" + req.TaskId
	geoServerURL := l.svcCtx.Config.GeoServerConfig.GeoServerURL
	workspace := l.svcCtx.Config.GeoServerConfig.Workspace

	username := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password
	storeURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s?recurse=true", geoServerURL, workspace, storeName)

	client := &http.Client{}
	_req, err := http.NewRequest("DELETE", storeURL, nil)
	if err != nil {
		return nil, err
	}

	_req.SetBasicAuth(username, password)
	resp, err := client.Do(_req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return &types.ErrorResponse{
			Code:    404,
			Message: "Resource not found",
		}, nil
	}

	return util.ParseErrorCode(err), nil
}
