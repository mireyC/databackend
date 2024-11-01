package types

import (
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
	"math"
	"strings"
)

type Line struct {
	Dx, Dy   float64        // 单位向量 (dx, dy)
	Index    int            // line 在数组中的 下标
	Points   [][]float64    // line 中coordinates 经纬度的信息
	Geometry Geometry       // 存到数据库的信息
	LineORB  orb.LineString // line中点的 orbLinestring
}

// ConvertToWKT
// 将 Geometry 转换为 WKT 格式的函数
// MULTILINESTRING 需要多条线段，每条线段用括号括起来，点之间用逗号分隔，线段之间也用逗号分隔
// MULTILINESTRING((116.391010 40.022661, 116.390711 40.022666), (116.390694 40.022667, 116.390422 40.022673, 116.390221 40.022675))
func ConvertToWKT(geometry Geometry) string {
	// 检查类型，当前只支持 MultiLineString
	if geometry.Type != "MultiLineString" {
		return ""
	}

	var wktBuilder strings.Builder

	// 开始构建 MULTILINESTRING 的 WKT 头部
	wktBuilder.WriteString("MULTILINESTRING(")

	// 遍历每个线段
	for i, line := range geometry.Coordinates {
		if i > 0 {
			wktBuilder.WriteString(", ")
		}
		wktBuilder.WriteString("(")
		// 遍历每条线段的点
		for j, point := range line {
			if j > 0 {
				wktBuilder.WriteString(", ")
			}
			// 将每个点的经纬度（和高度，如果有的话）写入 WKT
			wktBuilder.WriteString(fmt.Sprintf("%f %f", point[0], point[1]))
			if len(point) > 2 {
				wktBuilder.WriteString(fmt.Sprintf(" %f", point[2]))
			}
		}
		wktBuilder.WriteString(")")
	}

	wktBuilder.WriteString(")")

	return wktBuilder.String()
}

// ConvertPointsToLineString
// 将 [][]float64  Points 转换为 orb.LineString
func ConvertPointsToLineString(points [][]float64) orb.LineString {
	lineString := make(orb.LineString, len(points))
	for i, point := range points {
		if len(point) >= 2 {
			lineString[i] = orb.Point{point[0], point[1]} // 经度和纬度分别对应 x 和 y
		}
	}
	return lineString
}

// CreatePolygonByLines 使用两条线创建一个四边形（封闭区域）
func CreatePolygonByLines(line1, line2 orb.LineString) (orb.Polygon, bool) {
	// 使用 line1 起点和终点，和 line2 起点和终点构建四边形
	if len(line1) <= 1 || len(line2) <= 1 {
		return orb.Polygon{}, false
	}

	coords := []orb.Point{
		line1[0],
		line1[len(line1)-1],
		line2[len(line2)-1],
		line2[0],
		line1[0], // 闭合四边形
	}
	return orb.Polygon{coords}, true
}

// IsLineIntersectingPolygon 检查线是否与多边形相交
func IsLineIntersectingPolygon(polygon orb.Polygon, line orb.LineString) bool {
	for _, point := range line {
		if planar.RingContains(polygon[0], point) {
			return true
		}
	}
	return false
}

// MinDistanceBetweenLines
// 计算两条 LineString 之间的最小球面距离
func MinDistanceBetweenLines(line1, line2 orb.LineString) float64 {
	minDistance := math.MaxFloat64

	// 遍历 line1 和 line2 中的每个点对，计算球面距离
	for _, p1 := range line1 {
		for _, p2 := range line2 {
			distance := geo.Distance(p1, p2) // 使用球面距离计算
			if distance < minDistance {
				minDistance = distance
			}
		}
	}
	return minDistance
}

// CheckLimitedDeviation 检查较短的一条线段到较长线段的投影是否占较长线段长度的 70% 以上
func CheckLimitedDeviation(lineA, lineB Line) bool {
	// 计算每条线的长度
	lengthA := calculateLineLength(lineA.LineORB)
	lengthB := calculateLineLength(lineB.LineORB)

	// 找出较长和较短的线段及其长度
	var longerLine, shorterLine orb.LineString
	var longerLength float64
	//log.Println(shorterLength)
	if lengthA > lengthB {
		longerLine, shorterLine = lineA.LineORB, lineB.LineORB
		longerLength = lengthA
	} else {
		longerLine, shorterLine = lineB.LineORB, lineA.LineORB
		longerLength = lengthB
	}

	// 计算较短线段在较长线段上的投影长度
	projectionLength := calculateProjectionLength(longerLine, shorterLine)

	// 判断投影长度是否占较长线段长度的 70% 以上
	return projectionLength >= longerLength*0.7
}

// calculateProjectionLength 计算较短线段在较长线段上的投影长度
func calculateProjectionLength(longerLine, shorterLine orb.LineString) float64 {
	projectionLength := 0.0
	for i := 1; i < len(shorterLine); i++ {
		// 计算 shorterLine 上的每个线段在 longerLine 上的投影
		projectionLength += calculatePointProjection(shorterLine[i-1], shorterLine[i], longerLine)
	}
	return projectionLength
}

// calculatePointProjection 计算 shorterLine 上两个点之间线段在 longerLine 上的投影长度
func calculatePointProjection(p1, p2 orb.Point, longerLine orb.LineString) float64 {
	maxProjection := 0.0
	for j := 1; j < len(longerLine); j++ {
		segment := orb.LineString{longerLine[j-1], longerLine[j]}
		projection := calculateProjectionOnSegment(p1, p2, segment)
		if projection > maxProjection {
			maxProjection = projection
		}
	}
	return maxProjection
}

// calculateProjectionOnSegment 计算线段 (p1, p2) 在 longerLine 上某一段 segment 的投影长度
func calculateProjectionOnSegment(p1, p2 orb.Point, segment orb.LineString) float64 {
	// 计算 (p1, p2) 在 segment 上的投影长度
	segmentDirection := geo.Bearing(segment[0], segment[1])
	lineDirection := geo.Bearing(p1, p2)

	// 计算投影的比例
	cosAngle := math.Cos(segmentDirection - lineDirection)
	if cosAngle < 0 {
		return 0
	}

	// 计算投影长度
	projectionLength := geo.Distance(p1, p2) * cosAngle
	return projectionLength
}

// calculateLineLength 计算 Line 的长度
func calculateLineLength(line orb.LineString) float64 {
	length := 0.0
	for i := 1; i < len(line); i++ {
		length += geo.Distance(line[i-1], line[i])
	}
	return length
}

// CheckLimitedDeviation 检查两条线段的长度差是否在较长线段的 30% 范围内
//func CheckLimitedDeviation(lineA, lineB Line) bool {
//	// 计算每条线的长度
//	lengthA := calculateLineLength(lineA.LineORB)
//	lengthB := calculateLineLength(lineB.LineORB)
//
//	// 找出较长和较短的长度
//	longer, shorter := math.Max(lengthA, lengthB), math.Min(lengthA, lengthB)
//
//	// 判断长度差是否在较长线段的 30% 范围内
//	return math.Abs(longer-shorter) <= longer*0.3
//}

// 计算 Line 的长度
//func calculateLineLength(line orb.LineString) float64 {
//	length := 0.0
//	for i := 1; i < len(line); i++ {
//		length += geo.Distance(line[i-1], line[i])
//	}
//	return length
//}
