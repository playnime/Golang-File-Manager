package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AutoMigrate(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id           SERIAL PRIMARY KEY,
		login        VARCHAR(20) NOT NULL UNIQUE,
		password     VARCHAR(255) NOT NULL,
		created_at   TIMESTAMP DEFAULT NOW(),
		totp_secret  VARCHAR(64) NOT NULL DEFAULT '',
		totp_enabled BOOLEAN NOT NULL DEFAULT FALSE
	);

	CREATE TABLE IF NOT EXISTS files (
		id          SERIAL PRIMARY KEY,
		user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		file_name   VARCHAR(255) NOT NULL,
		path        TEXT NOT NULL,
		size_bytes  BIGINT,
		mime_type   VARCHAR(100),
		gd          INTEGER DEFAULT 0,
		uploaded_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id         SERIAL PRIMARY KEY,
		user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token      VARCHAR(255) NOT NULL UNIQUE,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		confirmed  BOOLEAN NOT NULL DEFAULT TRUE
	);
	`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("AutoMigrate: %w", err)
	}

	fmt.Println("migrations applied successfully")
	return nil
}
