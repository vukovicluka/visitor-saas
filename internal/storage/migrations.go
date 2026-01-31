package storage

import (
	"context"
	"fmt"
)

func (db *DB) migrate(ctx context.Context) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS page_views (
			id           BIGSERIAL PRIMARY KEY,
			domain       TEXT NOT NULL,
			path         TEXT NOT NULL,
			referrer     TEXT NOT NULL DEFAULT '',
			country_code TEXT NOT NULL DEFAULT '',
			visitor_hash TEXT NOT NULL DEFAULT '',
			created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE INDEX IF NOT EXISTS idx_page_views_domain_created
			ON page_views(domain, created_at)`,

		`CREATE INDEX IF NOT EXISTS idx_page_views_path
			ON page_views(domain, path, created_at)`,

		`CREATE INDEX IF NOT EXISTS idx_page_views_visitor
			ON page_views(domain, visitor_hash, created_at)`,

		`CREATE TABLE IF NOT EXISTS daily_salts (
			date DATE PRIMARY KEY,
			salt TEXT NOT NULL
		)`,
	}

	for _, m := range migrations {
		if _, err := db.pool.Exec(ctx, m); err != nil {
			return fmt.Errorf("exec migration: %w", err)
		}
	}

	return nil
}
