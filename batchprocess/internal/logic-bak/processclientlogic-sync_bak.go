package logic_bak

import (
	"batchprocess/internal/repository"
	"batchprocess/internal/svc"
	"batchprocess/internal/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"net/http"
	"os"
	"time"
)

type ProcessClientLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessClientLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessClientLogic {
	return &ProcessClientLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessClientLogic) ProcessClient(req *types.ProcessClientReq) (resp types.ErrorResponse, err error) {

	//laneChan := make(chan *repository.FeatureCollection)
	//markChan := make(chan *repository.FeatureCollection)
	ctx := context.Background()

	//go func() {
	s := repository.NewProcessClient(l.svcCtx.PgDb)
	p1, _ := s.SelectMarkLineByLimit(ctx, int(req.MarkOffset), int(req.MarkLimit))
	//markChan <- p11
	//}()

	//go func() {
	//	s := repository.NewProcessClient(l.svcCtx.PgDb)
	p2, _ := s.SelectLaneLineByLimit(ctx, int(req.LaneOffset), int(req.LaneLimit))
	//	laneChan <- p22
	//
	//}()

	//p1 := <-markChan
	//p2 := <-laneChan

	p1Json, _ := json.Marshal(p1)
	p2Json, _ := json.Marshal(p2)

	requestBody := map[string]interface{}{
		"had_lane_mark_line": json.RawMessage(p1Json), // 保持原始 JSON 格式
		"had_lane_link":      json.RawMessage(p2Json), // 保持原始 JSON 格式
		"mesh":               "1287290",               // 普通字符串
	}

	// 将 requestBody 序列化成 JSON
	finalJson, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshaling final JSON:", err)
		return
	}

	//fmt.Println("Final JSON:", string(finalJson))

	// 构造 POST 请求
	req_, err := http.NewRequest("POST", "http://127.0.0.1:8889/batchProcess", bytes.NewBuffer(finalJson))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req_.Header.Set("Content-Type", "application/json")

	// 记录开始时间
	startTime := time.Now()
	// 打开或创建一个日志文件，使用 os.O_APPEND 以追加方式写入日志
	logFile, err := os.OpenFile("processClient_log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// 创建两个 Logger 实例
	infoLogger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	// 发起请求

	//go func() {
	client := &http.Client{}
	resp_, err := client.Do(req_)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp_.Body.Close()
	//}()

	totalTime := time.Since(startTime)
	str := fmt.Sprintf("Processing Client completed, total process had_lane_mark_line count: %d,  total process had_lane_link: %d, total time costs %v", req.MarkLimit, req.LaneLimit, totalTime)
	infoLogger.Println(str)

	resp.Code = 200
	resp.Message = "success"
	return
}
