package serviceUsers

import (
	"lab2pepe/server/internal/users/repositoryUsers"
)

type ServiceUsers struct {
	RepoFiles *repositoryUsers.RepoUsers
}

func NewServiceUsers(repo *repositoryUsers.RepoUsers) *ServiceUsers {
	return &ServiceUsers{RepoFiles: repo}
}

func (s *ServiceUsers) AddUsers(login string, password string) error {

	return nil //s.RepoFiles.CreateUsers(login, password)
}
