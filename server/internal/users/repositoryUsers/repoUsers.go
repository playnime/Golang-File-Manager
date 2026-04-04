package repositoryUsers

import (
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

func (r *RepoUsers) CreateUsers(login string, password string) error {

	return nil
}
