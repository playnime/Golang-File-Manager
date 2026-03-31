package repositoryUsers

import "lab2pepe/server/internal/core/domain"

type RepoUsers struct {
	repo *domain.Repository
}

func NewHRepoFiles(repo *domain.Repository) *RepoUsers {
	return &RepoUsers{
		repo: repo,
	}
}

func (r *RepoUsers) CreateUsers() error {

	return nil
}
