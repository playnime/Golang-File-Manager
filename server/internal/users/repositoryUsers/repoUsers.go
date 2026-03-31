package repositoryUsers

import (
	"errors"
	"lab2pepe/server/internal/core/domain"
)

type RepoUsers struct {
	repo *domain.Repository
}

func NewHRepoFiles(repo *domain.Repository) *RepoUsers {
	return &RepoUsers{
		repo: repo,
	}
}

func (r *RepoUsers) CreateUsers(login string, password string) error {

	for _, v := range r.repo.Bd {
		if v.Login == login {
			return errors.New("Login exist's")
		}
	}
	r.repo.Bd = append(r.repo.Bd, domain.User{
		Login:    login,
		Password: password,
		Messages: nil,
	})
	return nil
}
