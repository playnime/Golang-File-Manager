package repositoryUsers

import (
	"context"
	"fmt"
	"lab2pepe/server/internal/core/md5MY"

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
