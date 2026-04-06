package serviceUsers

import (
	"context"
	"lab2pepe/server/internal/users/repositoryUsers"
)

type ServiceUsers struct {
	RepoFiles *repositoryUsers.RepoUsers
}

func NewServiceUsers(repo *repositoryUsers.RepoUsers) *ServiceUsers {
	return &ServiceUsers{RepoFiles: repo}
}

func (s *ServiceUsers) AddUsers(ctx context.Context, login string, password string) error {

	return s.RepoFiles.CreateUser(ctx, login, password) //s.RepoFiles.CreateUsers(login, password)
}
