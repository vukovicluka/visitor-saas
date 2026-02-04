package storage

import (
	"context"
	"fmt"
	"visitor/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func New(ctx context.Context, databaseUrl string) (*DB, error) {
	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	db := &DB{pool: pool}

	if err := db.migrate(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) InsertPageView(ctx context.Context, pv *model.PageView) error {
	_, err := db.pool.Exec(ctx,
		`INSERT INTO page_views (domain, path, referrer, country_code, visitor_hash)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (domain, path, visitor_hash, (created_at::date)) DO NOTHING`,
		pv.Domain, pv.Path, pv.Referrer, pv.CountryCode, pv.VisitorHash)

	return err
}
