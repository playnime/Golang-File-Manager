package serviceUsers

import (
	"context"
	"lab2pepe/server/internal/users/repositoryUsers"
)

type ServiceUsers struct {
	RepoUsers *repositoryUsers.RepoUsers
}

func NewServiceUsers(repo *repositoryUsers.RepoUsers) *ServiceUsers {
	return &ServiceUsers{RepoUsers: repo}
}

func (s *ServiceUsers) AddUsers(ctx context.Context, login string, password string) error {

	return s.RepoUsers.CreateUser(ctx, login, password) //s.RepoFiles.CreateUsers(login, password)
}

func (s *ServiceUsers) Login(ctx context.Context, login, password string) (string, error) {
	return s.RepoUsers.Login(ctx, login, password)
}
func (s *ServiceUsers) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	return s.RepoUsers.GetUserIDByToken(ctx, token)
}
