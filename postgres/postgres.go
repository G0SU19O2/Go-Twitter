package postgres

import (
	"context"
	"log"

	"github.com/G0SU19O2/Go-Twitter/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, conf *config.Config) *DB {
	dbConf, err := pgxpool.ParseConfig(conf.Database.URL)
	if err != nil {
		log.Fatalf("can't parse database config: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, dbConf)
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}
	db := &DB{Pool: pool}
	db.Pool.Ping(ctx)
	return db
}

func (db *DB) Ping(ctx context.Context) {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatalf("can't ping database: %v", err)
	}
	log.Println("Database connection is alive")
}

func (db *DB) Close() {
	db.Pool.Close()
}
