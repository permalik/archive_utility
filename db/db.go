package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"
)

var pool *sql.DB

func InitDB() *sql.DB {
	dsn := os.Getenv("DSN")
	var err error
	pool, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("unable to use dsn", err)
	}

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	return pool
}

func Ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database:\n%v", err)
	}
}
