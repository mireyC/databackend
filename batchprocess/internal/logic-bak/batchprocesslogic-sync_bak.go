package logic_bak

import (
	"batchprocess/internal/repository"
	"batchprocess/internal/svc"
	"batchprocess/internal/types"
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib" // 导入 pgx 驱动
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"math"
	"os"
	"time"
)

type BatchProcessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchProcessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchProcessLogic {
	return &BatchProcessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BatchProcess
// 处理 车道线（had_lane_mark_link） 生成中心线，批量插入数据库
func (l *BatchProcessLogic) BatchProcess(req *types.BatchProcessReq) (resp types.ErrorResponse, err error) {

	// var muIsGenarated sync.Mutex

	// 打开或创建一个日志文件，使用 os.O_APPEND 以追加方式写入日志
	var logFile, _ = os.OpenFile("batch_process_log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close() // 在程序结束时关闭文件

	// 创建两个 Logger 实例：一个用于 info，一个用于 error
	var infoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	//errorLogger := log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	//infoLogger.Println("isGerated: ", isGenarated)
	// 记录开始时间
	startTime := time.Now()

	var hadLandMarkLine types.HadLaneMarkLine
	hadLandMarkLine = req.Had_lane_mark_line

	var hadLaneLink types.HadLaneLink

	hadLaneLink = req.Had_lane_link
	//exeitsLines := make([]types.Line, 0)
	//infoLogger.Println("hadLaneLink: ", hadLaneLink)
	//infoLogger.Println("before deal exitsLines：", exeitsLines)
	// 处理已经存在的中心线
	exeitsLines := dealAlreadyExitsLines(hadLaneLink)
	infoLogger.Println("已经生成的中心线 数量：", len(exeitsLines))
	//infoLogger.Println(exeitsLines)
	//markLines := make([]types.Line, 0)
	var markLines []types.Line
	//var muMarkLines sync.Mutex
	//var wgMarkLines sync.WaitGroup

	for index, feature := range hadLandMarkLine.Features {
		//wgMarkLines.Add(1)
		//indexCopy := index
		//featureCopy := feature
		//go func(f types.MarkLinkFeature, index int) {
		//	defer wgMarkLines.Done()

		// 处理线段加入数组
		var dx, dy float64
		l := types.Line{}
		l.Index = index
		l.Geometry = feature.Geometry
		for _, line := range feature.Geometry.Coordinates {
			dx, dy = calculateDirection(line)
			l.Dx = dx
			l.Dy = dy
			l.Points = line
			l.LineORB = types.ConvertPointsToLineString(line)
		}
		//muMarkLines.Lock()
		markLines = append(markLines, l)
		//muMarkLines.Unlock()
		//infoLogger.Println("预处理marlines：", l.Index)
		//}(featureCopy, indexCopy)
	}

	//wgMarkLines.Wait()
	infoLogger.Println("要处理的车道线线数量: ", len(markLines))
	// 使用WaitGroup 确保 所有 中心线已经处理好
	//var wgCenterLine sync.WaitGroup
	var centerLins []types.Line
	var linesA []types.Line
	var linesB []types.Line
	//var muCenterLine sync.Mutex

	var isGenarated = make(map[int]bool)
	// 判断线段是否相邻，是否同向。如果相邻且同向则 生成 中心线
	//processCenterLineStartTime := time.Now()
	//go func() {
	processStartTime := time.Now()
	for _, line := range markLines {
		//infoLogger.Println("遍历所有中心线 ind: ", i)
		//wgCenterLine.Add(1)
		/*
			每个协程启动时会传递 line 变量作为参数，但是在使用闭包的并发场景中，这样传递参数会产生 闭包变量捕获问题。
			具体来说，line 是一个循环变量，它会在每次循环迭代时更新。
			如果直接在协程中使用它，会导致所有协程都共享同一个 line 变量的地址，进而导致并发问题和意外行为。

			将变量值在每次迭代时复制
			一种常见的解决方法是，在循环中显式复制 line 的值给一个局部变量，然后传递该局部变量到协程中。
			这样可以确保每个协程都使用不同的 line 实例。
		*/
		l := line
		//lCopy := line
		//go func(l types.Line) { // l -> 新来的线段，
		//infoLogger.Println("判断下标是否一致：", l.Index)
		//defer wgCenterLine.Done()

		//muMarkLines.Lock()
		//defer muMarkLines.Unlock()
		if isAlreadyGenarated(isGenarated, l.Index) {
			//return
			continue
		}
		minDistance := 100.0
		minDistLineIndex := -1

		dis := 0.0
		for ind, li := range markLines {
			if li.Index == l.Index {
				continue
			}
			if isAlreadyGenarated(isGenarated, li.Index) {
				continue
			}

			dis = types.MinDistanceBetweenLines(l.LineORB, li.LineORB)
			if dis < 10 && dis > 1 && dis < minDistance {
				minDistance = dis
				minDistLineIndex = ind
			}
		}

		if minDistLineIndex != -1 && checkSameDireton(markLines[minDistLineIndex], l) && types.CheckLimitedDeviation(markLines[minDistLineIndex], l) {

			//dis := types.MinDistanceBetweenLines(markLines[minDistLineIndex].LineORB, l.LineORB)
			//if dis > 1 {
			centerLin := genarateCenterLine(markLines[minDistLineIndex], l, isGenarated, infoLogger)

			if !checkAlreadyExits(centerLin, exeitsLines) {

				centerLins = append(centerLins, centerLin)
				linesA = append(linesA, markLines[minDistLineIndex])
				linesB = append(linesB, l)
				isGenarated[markLines[minDistLineIndex].Index] = true
				isGenarated[l.Index] = true
				//}
				//infoLogger.Println(isGenarated)
			}
			//else {
			//	infoLogger.Println("不匹配：", l.Index, "  和 ：", markLines[minDistLineIndex].Index)
			//}
		}
		//} else {
		//	infoLogger.Println("没有任何线与 该线匹配：", l.Index)
		//}

		//}(lCopy)
	}
	//}()

	//go func() {
	//	wgCenterLine.Wait() // 等待所有中心线处理好
	//infoLogger.Println("centerLine 数量： ", len(centerLins))
	//processCenterLineTotalTime := time.Since(processCenterLineStartTime)
	//str := fmt.Sprintf("Batch processing generate All center_line completed, total process mark_lines count: %d, total time costs %v", len(markLines), processCenterLineTotalTime)
	//infoLogger.Println(str)
	processTotalTime := time.Since(processStartTime)
	str := fmt.Sprintf("同步批处理，生成centLines %d, Total time：%v", len(centerLins), processTotalTime)
	infoLogger.Println(str)

	saveStartTime := time.Now()
	repo := repository.NewHadLaneCenterLineRepo(l.svcCtx.Db)
	ctx := context.Background()
	er := repo.SaveLAllToDatabase(ctx, centerLins, req.Mesh, linesA, linesB)
	if er != nil {
		log.Println("AAAFailed to insert to database: ", er)
	}
	er = repo.SaveAllToDatabase(ctx, centerLins, req.Mesh)
	if er != nil {
		log.Println("Failed to insert to database: ", er)
	}
	saveTotalTime := time.Since(saveStartTime)
	str = fmt.Sprintf("同步批处理：存入数据库 centerLine counts: %d, totalTime costs %v", len(centerLins), saveTotalTime)
	infoLogger.Println(str)
	//

	totalTime := time.Since(startTime)
	str = fmt.Sprintf("同步批处理：总时间 totalTime costs %v", totalTime)
	infoLogger.Println(str)
	infoLogger.Println()
	infoLogger.Println()
	infoLogger.Println()
	//	defer logFile.Close() // 在程序结束时关闭文件
	//}()
	resp.Code = 200
	resp.Message = "success"
	return
	//return &types.ErrorResponse{Code: 200, Message: "Success"}, nil
}

// isAlreadyGenarated
// 是否已经生成过中心线
func isAlreadyGenarated(m map[int]bool, index int) bool {
	//muIsGenarated.Lock()
	//defer muIsGenarated.Unlock()
	if m[index] {
		return true
	}
	return false
}

// checkAlreadyExits
// 中心线是否已经存在
func checkAlreadyExits(centerLine types.Line, exitsLines []types.Line) bool {

	length := len(exitsLines)
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			dis := types.MinDistanceBetweenLines(exitsLines[i].LineORB, exitsLines[j].LineORB)
			if dis < 10 {
				p, ok := types.CreatePolygonByLines(exitsLines[i].LineORB, exitsLines[j].LineORB)
				if !ok {
					continue
				}
				if types.IsLineIntersectingPolygon(p, centerLine.LineORB) {
					return true
				}
			}

		}
	}

	return false
}

// calculateDirection
// 计算车道的平均方向向量 (单位向量)
func calculateDirection(coordinates [][]float64) (float64, float64) {
	var totalX, totalY float64
	for i := 0; i < len(coordinates)-1; i++ {
		start := coordinates[i]
		end := coordinates[i+1]
		totalX += end[0] - start[0]
		totalY += end[1] - start[1]
	}

	magnitude := math.Sqrt(totalX*totalX + totalY*totalY)

	return totalX / magnitude, totalY / magnitude
}

// checkSameDireton
// 判断 lineA, lineB 是否 同向, 同向返回 true
// x1 * x2 + y1 * y2 => >0 同向, < 0 反向
func checkSameDireton(lineA, lineB types.Line) bool {

	x1 := lineA.Dx
	y1 := lineA.Dy
	x2 := lineB.Dx
	y2 := lineB.Dy

	return x1*x2+y1*y2 > 0
}

// caculateCenterPoint
// 计算两点间的中间点
// (x1, y1), (x2, y2) => ((x1 + x2) / 2, (y1 + y2) / 2)
func caculateCenterPoint(pointA, pointB []float64) []float64 {
	x1 := pointA[0]
	y1 := pointA[1]
	x2 := pointB[0]
	y2 := pointB[1]

	var res []float64
	res = append(res, (x1+x2)/2)
	res = append(res, (y1+y2)/2)
	return res
}

// genarateCenterLine
// 计算俩条线的中心线
func genarateCenterLine(lineA, lineB types.Line, isGenarated map[int]bool, infoLogger *log.Logger) types.Line {

	var points [][]float64

	LenA := len(lineA.Points)
	LenB := len(lineB.Points)
	if LenA < LenB {
		for i := 0; i < LenA; i++ {
			pointA := lineA.Points[i]
			pointB := lineB.Points[i]
			points = append(points, caculateCenterPoint(pointA, pointB))
			if i+1 < LenB {
				points = append(points, caculateCenterPoint(pointA, lineB.Points[i+1]))
			}
		}

		points = append(points, caculateCenterPoint(lineA.Points[LenA-1], lineB.Points[LenB-1]))
	} else {
		for i := 0; i < LenB; i++ {
			pointB := lineB.Points[i]
			pointA := lineA.Points[i]
			points = append(points, caculateCenterPoint(pointB, pointA))
			if i+1 < LenA {
				points = append(points, caculateCenterPoint(pointB, lineA.Points[i+1]))
			}
		}

		points = append(points, caculateCenterPoint(lineB.Points[LenB-1], lineA.Points[LenA-1]))
	}

	l := types.Line{}
	l.Geometry.Type = lineA.Geometry.Type
	l.Points = points
	l.Geometry.Coordinates = append(l.Geometry.Coordinates, points)

	return l
}

// dealAlreadyExitsLines
// 处理已经存在在中心线 转为 []Line
func dealAlreadyExitsLines(geoJson types.HadLaneLink) []types.Line {

	var lines []types.Line

	for index, feature := range geoJson.Features {
		// 处理线段加入数组
		var dx, dy float64
		l := types.Line{}
		l.Index = index
		l.Geometry = feature.Geometry
		for _, line := range feature.Geometry.Coordinates {
			dx, dy = calculateDirection(line)
			l.Dx = dx
			l.Dy = dy
			l.Points = line
			l.LineORB = types.ConvertPointsToLineString(line)
			lines = append(lines, l)
		}
	}

	return lines
}
