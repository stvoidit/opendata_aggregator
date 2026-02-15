// Package store - имплиментация работы с БД
package store

import (
	"context"
	"fmt"
	"opendataaggregator/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB - ...
type DB struct {
	pool *pgxpool.Pool
}

// NewDB - ...
func NewDB(cnf *config.Config, workersCount int) (*DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable application_name=%s",
		cnf.DB.Host,
		cnf.DB.Port,
		cnf.DB.Login,
		cnf.DB.Password,
		cnf.DB.DBName,
		cnf.DB.DBName)
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	poolConfig.MinConns = int32(workersCount / 2)
	poolConfig.MaxConns = int32(workersCount)
	poolConfig.MaxConnLifetime = time.Minute * 10
	poolConfig.MaxConnIdleTime = time.Minute * 5
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &DB{pool: pool}, nil
}

// Close - закрытие соединений
func (db *DB) Close() { db.pool.Close() }

func isNullStr(s string) *string {
	if len(s) == 0 {
		return nil
	}
	return &s
}

// func isStrNull(s *string) string {
// 	if s == nil {
// 		return ""
// 	}
// 	return *s
// }
