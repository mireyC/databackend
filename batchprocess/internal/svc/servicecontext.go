package svc

import (
	"batchprocess/internal/config"
	"context"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib" // 导入 pgx 驱动
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"log"
)

func NewPgxConn(dbSource string) *pgx.Conn {
	// 建立连接
	connString := dbSource

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

type ServiceContext struct {
	Config config.Config
	Db     sqlx.SqlConn
	PgDb   *pgx.Conn
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Db:     sqlx.NewSqlConn("pgx", c.DataSource),
		PgDb:   NewPgxConn(c.DataSource),
	}
}
