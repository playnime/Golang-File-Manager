package serviceUsers

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"fmt"
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
	userID, err := s.RepoUsers.GetUserByCredentials(ctx, login, password)
	if err != nil {
		return "", err
	}

	b := make([]byte, 32)
	if _, err := crypto_rand.Read(b); err != nil {
		return "", fmt.Errorf("serviceUsers.Login: generate token: %w", err)
	}
	token := hex.EncodeToString(b)

	if err := s.RepoUsers.SaveSession(ctx, userID, token); err != nil {
		return "", err
	}

	fmt.Println("User authorized:", login)
	return token, nil
}

func (s *ServiceUsers) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	return s.RepoUsers.GetUserIDByToken(ctx, token)
}
