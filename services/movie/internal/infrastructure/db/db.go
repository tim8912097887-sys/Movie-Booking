package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/shutdown"
)

type DbConfig struct {
   Ctx context.Context
   Logger *slog.Logger
   DbUrl string
   ShutdownManager *shutdown.ShutdownManager
}

type Db struct {
	dbConfig DbConfig
}

func NewDb(dbConfig DbConfig) *Db {
	return &Db{dbConfig: dbConfig}
}

func (d *Db) InitDB() (*pgxpool.Pool,error) {

	dbpool, err := pgxpool.New(d.dbConfig.Ctx, d.dbConfig.DbUrl)
	if err != nil {
		return nil, err
	}
	 if err := dbpool.Ping(d.dbConfig.Ctx); err != nil {
        return nil,err
    }

	d.dbConfig.Logger.Info("connected to database",slog.String("url", d.dbConfig.DbUrl))

	d.dbConfig.ShutdownManager.Register(d.CloseDB(dbpool))
	d.dbConfig.Logger.Info("registered shutdown handler for db")
	return dbpool, nil
}

func (d *Db) CloseDB(dbpool *pgxpool.Pool) func (ctx context.Context) error {
	return func(ctx context.Context) error {
		dbpool.Close()
		return nil
	}
}