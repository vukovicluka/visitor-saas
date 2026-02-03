package dashboard

import (
	"context"
	"fmt"
	"visitor/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	pool *pgxpool.Pool
}

func NewQueries(pool *pgxpool.Pool) *Queries {
	return &Queries{pool: pool}
}

func (q *Queries) Summary(ctx context.Context, domain string, days int) (*model.SummaryStats, error) {
	stats := &model.SummaryStats{}

	err := q.pool.QueryRow(ctx, `SELECT COUNT(*), COUNT(DISTINCT visitor_hash)
								FROM page_views
								WHERE domain = $1 AND created_at >= NOW() - make_interval(days => $2)`,
							domain, days).Scan(&stats.TotalViews, &stats.UniqueVisitors)
	if err != nil {
		return nil, fmt.Errorf("summary totals: %w", err)
	}

	rows, err := q.pool.Query(ctx, 
							`SELECT created_at::date::text AS date,
							COUNT(*) AS views,
							COUNT(DISTINCT visitor_hash) AS visitors
							FROM page_views
							WHERE domain = $1 AND created_at >= NOW() - make_interval(days => $2)
							GROUP BY created_at::date
							ORDER BY created_at::date`, domain, days)
	if err != nil {
		return nil, fmt.Errorf("daily stats: %w", err)
	}	

	defer rows.Close()

	for rows.Next() {
		var d model.DailyStat
		if err := rows.Scan(&d.Date, &d.Views, &d.Visitors); err != nil {
			return nil, fmt.Errorf("scan daily stat: %w", err)
		}
		stats.ViewsPerDay = append(stats.ViewsPerDay, d)
	}

	return stats, rows.Err()
}

func (q *Queries) Pages(ctx context.Context, domain string, days int) ([]model.PageStats, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT path, COUNT(*) AS views, COUNT(DISTINCT visitor_hash) AS visitors
		 FROM page_views
		 WHERE domain = $1 AND created_at >= NOW() - make_interval(days => $2)
		 GROUP BY path
		 ORDER BY views DESC
		 LIMIT 20`,
		domain, days)
	if err != nil {
		return nil, fmt.Errorf("top pages: %w", err)
	}
	defer rows.Close()

	var pages []model.PageStats
	for rows.Next() {
		var p model.PageStats
		if err := rows.Scan(&p.Path, &p.Views, &p.Visitors); err != nil {
			return nil, fmt.Errorf("scan page: %w", err)
		}
		pages = append(pages, p)
	}

	return pages, rows.Err()
}

func (q *Queries) Referrers(ctx context.Context, domain string, days int) ([]model.ReferrerStats, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT referrer, COUNT(*) AS views, COUNT(DISTINCT visitor_hash) AS visitors
		 FROM page_views
		 WHERE domain = $1 AND referrer != '' AND created_at >= NOW() - make_interval(days => $2)
		 GROUP BY referrer
		 ORDER BY views DESC
		 LIMIT 20`,
		domain, days)
	if err != nil {
		return nil, fmt.Errorf("top referrers: %w", err)
	}
	defer rows.Close()

	var refs []model.ReferrerStats
	for rows.Next() {
		var r model.ReferrerStats
		if err := rows.Scan(&r.Referrer, &r.Views, &r.Visitors); err != nil {
			return nil, fmt.Errorf("scan referrer: %w", err)
		}
		refs = append(refs, r)
	}

	return refs, rows.Err()
}