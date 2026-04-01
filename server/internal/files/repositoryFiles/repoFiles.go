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

func (r *RepoFiles) RepoAddFile(login string, password string, message string) error {
	var k int = -1
	for i, v := range r.repo.Bd {
		if login == v.Login {
			if password == v.Password {
				k = i
			} else {
				return errors.New("password is not correct")
			}
			break
		}
	}

	if k != -1 {
		r.repo.Bd[k].Messages = append(r.repo.Bd[k].Messages, message)
	} else {
		return errors.New("Login is not correct")
	}

	return nil
}

func (r *RepoFiles) RepoGetFiles(login string, password string) ([]string, error) {
	var k int = -1
	for i, v := range r.repo.Bd {
		if login == v.Login {
			if password == v.Password {
				k = i
			} else {
				return nil, errors.New("password is not correct")
			}
			break
		}
	}

	if k != -1 {
		return r.repo.Bd[k].Messages, nil
	}

	return nil, errors.New("Login is not correct")
}

//еще будет функция скачивания файла
