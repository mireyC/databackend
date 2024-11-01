package repository

import (
	"batchprocess/internal/types"
	"context"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib" // 导入 pgx 驱动
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"log"
	"strings"
	"time"
)

type HadLaneCenterLine struct {
	Gid       int64   `json:"gid" db:"gid"` // 主键
	ID        float64 `json:"id" db:"id"`
	LaneLink  string  `json:"lane_link_" db:"lane_link_"`
	LgID      float64 `json:"lg_id" db:"lg_id"`
	SeqNum    float64 `json:"seq_num" db:"seq_num"`
	CreateTim float64 `json:"create_tim" db:"create_tim"`
	UpdateTim float64 `json:"update_tim" db:"update_tim"`
	MaterialI string  `json:"material_i" db:"material_i"`
	Operator  string  `json:"operator" db:"operator"`
	Version   float64 `json:"version" db:"version"`
	Mesh      string  `json:"mesh" db:"mesh"`
	LaneType  float64 `json:"lane_type" db:"lane_type"`
	IsInterse float64 `json:"is_interse" db:"is_interse"`
	Geom      string  `json:"geom" db:"geom"`
}

type HadLaneCenterLineRepository interface {
	SaveAllToDataBase(ctx context.Context, centerLines []types.Line, mesh string) error
}

type hadLaneCenterLineRepo struct {
	db sqlx.SqlConn
}

func NewHadLaneCenterLineRepo(db sqlx.SqlConn) hadLaneCenterLineRepo {
	return hadLaneCenterLineRepo{
		db: db,
	}
}

func (h *hadLaneCenterLineRepo) SaveLAllToDatabase(ctx context.Context, centerLines []types.Line, mesh string, linesA []types.Line, linesB []types.Line) error {

	// 构建批量插入 SQL
	var values []interface{}
	var placeholders []string

	now := time.Now().Unix()
	//log.Println("centerLines count: ", len(centerLines))
	for i := 0; i < len(centerLines); i++ {
		centerLineWKT := types.ConvertToWKT(types.Geometry(centerLines[i].Geometry))
		la := types.ConvertToWKT(types.Geometry(linesA[i].Geometry))
		lb := types.ConvertToWKT(types.Geometry(linesB[i].Geometry))
		createTime := now
		updateTime := now

		// 构建占位符，i 是行数，构建 4个字段 (createTime, updateTime, mesh, geom)
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, ST_GeomFromText($%d), ST_GeomFromText($%d), ST_GeomFromText($%d))", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))

		// 加入对应的值：createTime, updateTime, mesh 和 WKT 几何数据
		values = append(values, createTime, updateTime, mesh, centerLineWKT, la, lb)
	}

	// 如果没有有效数据则返回错误
	if len(placeholders) == 0 {
		return errors.New("没有可出保存数据")
	}

	// 构建 SQL 语句
	query := fmt.Sprintf("INSERT INTO public.new_had_lane_center_line (create_tim, update_tim, mesh, geom, linea_geom, lineb_geom) VALUES %s", strings.Join(placeholders, ","))

	// 执行批量插入
	_, err := h.db.ExecCtx(ctx, query, values...)
	if err != nil {
		log.Println("Failed to batch insert into database:", err)
		return err
	}

	fmt.Println("Data batch inserted successfully")
	return nil
}

func (h *hadLaneCenterLineRepo) SaveAllToDatabase(ctx context.Context, centerLines []types.Line, mesh string) error {

	// 构建批量插入 SQL
	var values []interface{}
	var placeholders []string

	now := time.Now().Unix()
	//log.Println("centerLines count: ", len(centerLines))
	for i := 0; i < len(centerLines); i++ {
		centerLineWKT := types.ConvertToWKT(types.Geometry(centerLines[i].Geometry))
		//la := linesA[i]
		createTime := now
		updateTime := now

		// 构建占位符，i 是行数，构建 4个字段 (createTime, updateTime, mesh, geom)
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, ST_GeomFromText($%d))", i*4+1, i*4+2, i*4+3, i*4+4))

		// 加入对应的值：createTime, updateTime, mesh 和 WKT 几何数据
		values = append(values, createTime, updateTime, mesh, centerLineWKT)
	}

	// 如果没有有效数据则返回错误
	if len(placeholders) == 0 {
		return errors.New("没有可出保存数据")
	}

	// 构建 SQL 语句
	query := fmt.Sprintf("INSERT INTO public.had_lane_center_line (create_tim, update_tim, mesh, geom) VALUES %s", strings.Join(placeholders, ","))

	// 执行批量插入
	_, err := h.db.ExecCtx(ctx, query, values...)
	if err != nil {
		log.Println("Failed to batch insert into database:", err)
		return err
	}

	fmt.Println("Data batch inserted successfully")
	return nil
}
