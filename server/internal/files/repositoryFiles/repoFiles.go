package repositoryfiles

import (
	"errors"
	"lab2pepe/server/internal/core/domain"
)

type RepoFiles struct {
	repo *domain.Repository
}

func NewHRepoFiles(repo *domain.Repository) *RepoFiles {
	return &RepoFiles{
		repo: repo,
	}
}

func (r *RepoFiles) CreateMessage(login string, password string, message string) error {
	var k int = -1
	for i, v := range r.repo.Bd {
		if login == v.Login {
			if password == v.Password {
				k = i
				break
			}
		}
	}

	if k != -1 {
		r.repo.Bd[k].Messages = append(r.repo.Bd[k].Messages, message)
	} else {
		return errors.New("Login or password uncorrect")
	}

	return nil
}

func (r *RepoFiles) GetMessages(login string, password string) []string {
	var k int = -1
	for i, v := range r.repo.Bd {
		if login == v.Login {
			if password == v.Password {
				k = i
				break
			}
		}
	}

	if k != -1 {
		return r.repo.Bd[k].Messages
	}
	return nil
}
