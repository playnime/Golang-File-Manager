package repositoryUsers

import (
	"context"
	"encoding/hex"
	"fmt"
	"lab2pepe/server/internal/core/md5MY"

	crypto_rand "crypto/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoUsers struct {
	Pool *pgxpool.Pool
}

func NewRepoUsers(pool *pgxpool.Pool) *RepoUsers {
	return &RepoUsers{
		Pool: pool,
	}
}

func (r *RepoUsers) CreateUser(ctx context.Context, login, password string) error {
	hash := md5MY.HashToString(md5MY.HASH([]byte(password)))

	_, err := r.Pool.Exec(ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2)`,
		login, hash,
	)

	if err != nil {
		return fmt.Errorf("repoUsers.CreateUser: %w", err)
	}
	return nil
}

func (r *RepoUsers) Login(ctx context.Context, login, password string) (string, error) {
	hash := md5MY.HashToString(md5MY.HASH([]byte(password)))

	var userID int
	err := r.Pool.QueryRow(ctx,
		`SELECT id FROM users WHERE login = $1 AND password = $2`,
		login, hash,
	).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("repoUsers.Login: %w", err)
	}

	// генерируем токен
	b := make([]byte, 32)
	crypto_rand.Read(b)
	token := hex.EncodeToString(b)

	_, err = r.Pool.Exec(ctx,
		`INSERT 
		INTO sessions (user_id, token, expires_at)
        VALUES ($1, $2, NOW() + INTERVAL '24 hours')`,
		userID,
		token,
	)
	if err != nil {
		return "", fmt.Errorf("repoUsers.Login: save session: %w", err)
	}

	return token, nil
}

func (r *RepoUsers) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	var userID int
	err := r.Pool.QueryRow(ctx,
		`SELECT user_id FROM sessions 
		WHERE token = $1 AND expires_at > NOW()`,
		token,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("repoUsers.GetUserIDByToken: %w", err)
	}
	return userID, nil
}
