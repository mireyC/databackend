package geoContainer

import (
	"bytes"
	"context"
	"fmt"
	"geoserver/api/internal/util"
	"log"
	"net/http"
	"strings"
	"time"

	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDataSourceByPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDataSourceByPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDataSourceByPostLogic {
	return &CreateDataSourceByPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDataSourceByPostLogic) CreateDataSourceByPost(req *types.CreateDataSourceReq) (*types.ErrorResponse, error) {
	geoServerURL := l.svcCtx.Config.GeoServerConfig.GeoServerURL
	username := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password
	workSpace := l.svcCtx.Config.GeoServerConfig.Workspace
	storeType := l.svcCtx.Config.GeoServerConfig.StoreType
	if len(req.BucketUrl) > 0 && req.BucketUrl[0] == '/' {
		req.BucketUrl = req.BucketUrl[1:] // 去掉第一个字符 '/'
	}
	fileURL := "file:" + req.BucketUrl

	coverageStoreName := "bev" + "_" + req.TaskId
	// Check if workspace exists
	checkWorkspaceURL := fmt.Sprintf("%s/rest/workspaces/%s", geoServerURL, workSpace)
	workspaceReq, err := http.NewRequest("GET", checkWorkspaceURL, nil)
	if err != nil {
		fmt.Println("Error checking workspace: ", err)
		return nil, err
	}
	workspaceReq.SetBasicAuth(username, password)
	workspaceClient := &http.Client{}
	workspaceResp, err := workspaceClient.Do(workspaceReq)
	if err != nil {
		fmt.Println("Error checking workspace existence: ", err)
		return nil, err
	}
	defer workspaceResp.Body.Close()

	if workspaceResp.StatusCode == http.StatusNotFound {
		return &types.ErrorResponse{
			Code:    404,
			Message: fmt.Sprintf("Workspace does not exist, please create workspace %s first", workSpace),
		}, nil
	}

	// check if store already exists
	checkStoreURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s", geoServerURL, workSpace, coverageStoreName)
	checkReq, err := http.NewRequest("GET", checkStoreURL, nil)
	if err != nil {
		fmt.Println("err check store, ", err)
		return nil, err
	}

	checkReq.SetBasicAuth(username, password)
	client := &http.Client{}
	checkResp, err := client.Do(checkReq)
	if err != nil {
		fmt.Println("Error checking coverage store: ", err)
		return nil, err
	}

	defer checkResp.Body.Close()

	if checkResp.StatusCode == http.StatusOK {
		return &types.ErrorResponse{
			Code:    500,
			Message: "Data Source already exists",
		}, nil
	}

	// XML data for the coverage store
	coverageStoreXML := fmt.Sprintf(`<coverageStore>
  			<name>%s</name>
  			<type>%s</type>
  			<enabled>true</enabled>
  			<workspace>%s</workspace>
  			<url>%s</url>
	</coverageStore>`, coverageStoreName, storeType, workSpace, fileURL)

	// Create a new HTTP request
	createImageMosaicStoreURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores", geoServerURL, workSpace)
	_req, err := http.NewRequest("POST", createImageMosaicStoreURL, strings.NewReader(coverageStoreXML))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil, err
	}

	// Set headers and authentication
	_req.Header.Set("Content-Type", "text/xml")
	_req.SetBasicAuth(username, password)

	// Send the request
	client = &http.Client{}
	resp, err := client.Do(_req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 调用创建服务， 创建 ImageMosaic
	go func() {
		time.Sleep(1 * time.Second)
		imageName := coverageStoreName
		imageTitle := coverageStoreName
		storeName := coverageStoreName

		imageMosaicURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s/coverages", geoServerURL, workSpace, storeName)
		imageMosaicXML := fmt.Sprintf(`<coverage>
			<name>%s</name>
			<title>%s</title>
			<abstract>This is a sample image mosaic</abstract>
			<enabled>true</enabled>
			<store class="coverageStore">
				<name>%s</name>
			</store>
		</coverage>`, imageName, imageTitle, storeName)
		_req, er := http.NewRequest("POST", imageMosaicURL, bytes.NewBufferString(imageMosaicXML))
		if er != nil {
			log.Println("err create imageMosaic request during create store: , ", er)
		}

		_req.Header.Set("Content-Type", "text/xml")
		_req.SetBasicAuth(username, password)
		client = &http.Client{}
		resp, er := client.Do(_req)
		if er != nil {
			log.Println("err create imageMosaic during create store: , ", er)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {

			log.Println("err create imageMosaic during create store， code: , err,: , ", resp.StatusCode, er)

		}

	}()

	return util.ParseErrorCode(err), nil
}
