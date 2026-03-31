package repositoryfiles

import "lab2pepe/server/internal/core/domain"

type RepoFiles struct {
	repo *domain.Repository
}

func NewHRepoFiles(repo *domain.Repository) *RepoFiles {
	return &RepoFiles{
		repo: repo,
	}
}

func (r *RepoFiles) CreateFile() error {

	return nil
}
