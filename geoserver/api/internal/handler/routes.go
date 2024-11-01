// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package handler

import (
	"net/http"

	geoContainer "geoserver/api/internal/handler/geoContainer"
	"geoserver/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/createDataSource",
				Handler: geoContainer.CreateDataSourceByGetHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/createDataSource",
				Handler: geoContainer.CreateDataSourceByPostHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/createImageMosaicByStore",
				Handler: geoContainer.CreateImageMosaicByStoreHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/createImageMosaicByStoreV2",
				Handler: geoContainer.CreateImageMosaicByStoreV2Handler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/createImageMosaicByZip",
				Handler: geoContainer.CreateImageMosaicByZipHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/createImageMosaicByZipV2",
				Handler: geoContainer.CreateImageMosaicByZipV2Handler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/createImageMosaicStore",
				Handler: geoContainer.CreateImageMosaicStoreHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/del",
				Handler: geoContainer.DelGeoContainerHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/delDataSource",
				Handler: geoContainer.DelDataSourceByGetHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/delDataSource",
				Handler: geoContainer.DelDataSourceHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/delImageMosaic",
				Handler: geoContainer.DelImageMosaicHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/delImageMosaicStore",
				Handler: geoContainer.DelImageMosaicStoreHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/delImageMosaicStoreV2",
				Handler: geoContainer.DelImageMosaicStoreV2Handler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/delImageMosaicV2",
				Handler: geoContainer.DelImageMosaicV2Handler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/getServerStatus",
				Handler: geoContainer.GetServerStatusHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/start",
				Handler: geoContainer.StartGeoContainerHandler(serverCtx),
			},
		},
		rest.WithPrefix("/dataCenter/mdbServer/v1"),
	)
}