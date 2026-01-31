package hash

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
	pool *pgxpool.Pool
}

func NewManager(pool *pgxpool.Pool) *Manager {
	return &Manager{pool: pool}
}

func (m *Manager) GetHash(ctx context.Context, domain, ip, userAgent string) (string, error) {
	today := time.Now().UTC().Format("2006-01-02")

	salt, err := m.getSalt(ctx, today)
	if err != nil {
		return "", fmt.Errorf("get salt: %w", err)
	}

	raw := salt + ":" + domain + ":" + ip + ":" + userAgent
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:]), nil
}

func (m *Manager) getSalt(ctx context.Context, date string) (string, error) {
	var salt string
	err := m.pool.QueryRow(ctx, "SELECT salt from daily_salts WHERE date = $1", date).Scan(&salt)
	if err != nil {
		return salt, nil
	}

	// Generate new salt for today
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate random salt: %w", err)
	}
	salt = hex.EncodeToString(bytes)

	_, err = m.pool.Exec(ctx, 
			`INSERT INTO daily_salts (date, salt) 
			VALUES ($1, $2) ON CONFLICT (date) DO NOTHING`,
			date, salt)

	if err != nil {
		return "", fmt.Errorf("insert salt: %w", err)
	}		

	return salt, nil
}

func (m *Manager) cleanOldSalts(ctx context.Context) error {
	_, err := m.pool.Exec(ctx, "DELETE FROM daily_salts WHERE date < NOW() - INTERVAL  '2 days'")
	return err
}