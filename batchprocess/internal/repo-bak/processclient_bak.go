package repo_bak

//
//import (
//	"context"
//	"encoding/json"
//	"github.com/jackc/pgx/v4"
//)
//
//type HadLaneLink struct {
//	GID       int     `json:"gid"`        // integer
//	ID        float64 `json:"id"`         // numeric
//	LaneLink  string  `json:"lane_link"`  // character varying (varchar)
//	LGID      float64 `json:"lg_id"`      // numeric
//	SeqNum    float64 `json:"seq_num"`    // double precision
//	CreateTim float64 `json:"create_tim"` // numeric
//	UpdateTim float64 `json:"update_tim"` // numeric
//	MaterialI string  `json:"material_i"` // character varying (varchar)
//	Operator  string  `json:"operator"`   // character varying (varchar)
//	Version   float64 `json:"version"`    // double precision
//	Mesh      float64 `json:"mesh"`       // double precision
//	LaneType  float64 `json:"lane_type"`  // double precision
//	IsInterse float64 `json:"is_interse"` // double precision
//	Geom      string  `json:"geom"`       // USER-DEFINED (geometry)
//}
//
//type HadLaneMarkLink struct {
//	GID       int     `json:"gid"`        // integer
//	ID        float64 `json:"id"`         // numeric
//	LaneMark  string  `json:"lane_mark"`  // character varying (varchar)
//	CreateTim float64 `json:"create_tim"` // numeric
//	UpdateTim float64 `json:"update_tim"` // numeric
//	MaterialI string  `json:"material_i"` // character varying (varchar)
//	Operator  string  `json:"operator"`   // character varying (varchar)
//	Version   float64 `json:"version"`    // double precision
//	Mesh      string  `json:"mesh"`       // character varying (varchar)
//	MarkType  float64 `json:"mark_type"`  // double precision
//	MarkColor float64 `json:"mark_color"` // double precision
//	MarkWidth float64 `json:"mark_width"` // double precision
//	LeftLane  float64 `json:"left_lane"`  // numeric
//	RightLane float64 `json:"right_lane"` // numeric
//	Geom      string  `json:"geom"`       // USER-DEFINED (geometry)
//}
//
//type FeatureCollection struct {
//	Type     string    `json:"type"`
//	Name     string    `json:"name"`
//	Features []Feature `json:"features"`
//}
//
//type Feature struct {
//	Type       string                 `json:"type"`
//	Properties map[string]interface{} `json:"properties"`
//	Geometry   json.RawMessage        `json:"geometry"`
//}
//
//type ProcessClientRepository interface {
//	SelectMarkLineByLimit(ctx context.Context, offset int, limit int) (*FeatureCollection, error)
//	SelectLaneLineByLimit(ctx context.Context, offset int, limit int) (*FeatureCollection, error)
//}
//
//type ProcessClient struct {
//	db *pgx.Conn
//}
//
//func NewProcessClient(db *pgx.Conn) ProcessClientRepository {
//	return &ProcessClient{
//		db: db,
//	}
//}
//
//// SelectMarkLineByLimit
//// 车道线（had_lane_mark_link）
//func (p ProcessClient) SelectMarkLineByLimit(ctx context.Context, offset int, limit int) (*FeatureCollection, error) {
//	query := `
//        SELECT gid, id, ST_AsGeoJSON(geom) AS geom
//        FROM had_lane_mark_link
//        ORDER BY gid
//        LIMIT $1 OFFSET $2
//    `
//
//	rows, err := p.db.Query(ctx, query, limit, offset)
//	if err != nil {
//		return &FeatureCollection{}, err
//	}
//	defer rows.Close()
//
//	var features []Feature
//	for rows.Next() {
//		var geom json.RawMessage
//		var gid int
//		var id float64
//
//		if err := rows.Scan(&gid, &id, &geom); err != nil {
//			return &FeatureCollection{}, err
//		}
//
//		// 构造一个 Properties
//		properties := map[string]interface{}{
//			"gid":        gid,
//			"id":         id,
//			"lane_mark_": "k", // 可以为空
//			"create_tim": 1723086821509.0,
//			"update_tim": 1723086821509.0,
//			"material_i": "k", // 可以为空
//			"operator":   "auto",
//			"version":    0.0,
//			"mesh":       "1287290",
//			"mark_type":  1.0,
//			"mark_color": 1.0,
//			"mark_width": 1.2, // 可以为空
//			"left_lane_": 1314212940.0,
//			"right_lane": 1314212940.0,
//		}
//
//		feature := Feature{
//			Type:       "Feature",
//			Properties: properties,
//			Geometry:   geom,
//		}
//
//		features = append(features, feature)
//	}
//
//	if err := rows.Err(); err != nil {
//		return nil, err
//	}
//
//	featureCollection := &FeatureCollection{
//		Type:     "FeatureCollection",
//		Name:     "had_lane_mark_line",
//		Features: features,
//	}
//
//	return featureCollection, nil
//}
//
//func (p ProcessClient) SelectLaneLineByLimit(ctx context.Context, offset int, limit int) (*FeatureCollection, error) {
//
//	query := `
//        SELECT ST_AsGeoJSON(geom) AS geom
//        FROM had_lane_link
//        ORDER BY gid
//        LIMIT $1 OFFSET $2
//    `
//
//	rows, err := p.db.Query(context.Background(), query, limit, offset)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var features []Feature
//	for rows.Next() {
//		var geom json.RawMessage
//		if err := rows.Scan(&geom); err != nil {
//			return nil, err
//		}
//
//		// 创建一个 Feature 并填充空的 properties
//		feature := Feature{
//			Type: "Feature",
//			Properties: map[string]interface{}{
//				"gid":        1,
//				"id":         1,
//				"lane_link_": "kk",
//				"lg_id":      1,
//				"seq_num":    1.0,
//				"create_tim": 1.0,
//				"update_tim": 1.0,
//				"material_i": "kk",
//				"operator":   "k",
//				"version":    0.0,
//				"mesh":       "k",
//				"lane_type":  0.0,
//				"is_interse": 0.0,
//			},
//			Geometry: geom,
//		}
//
//		features = append(features, feature)
//	}
//
//	if err := rows.Err(); err != nil {
//		return nil, err
//	}
//
//	// 创建 FeatureCollection 并填充
//	featureCollection := &FeatureCollection{
//		Type:     "FeatureCollection",
//		Name:     "b",
//		Features: features,
//	}
//	return featureCollection, nil
//}
