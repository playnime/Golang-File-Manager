package serviceUsers

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"fmt"
	"lab2pepe/server/internal/users/repositoryUsers"

	"github.com/pquerna/otp/totp"
)

type ServiceUsers struct {
	RepoUsers *repositoryUsers.RepoUsers
}

func NewServiceUsers(repo *repositoryUsers.RepoUsers) *ServiceUsers {
	return &ServiceUsers{RepoUsers: repo}
}

type LoginResult struct {
	Token      string
	TmpToken   string
	Require2FA bool
}

func (s *ServiceUsers) AddUsers(ctx context.Context, login string, password string) error {
	return s.RepoUsers.CreateUser(ctx, login, password)
}

func (s *ServiceUsers) Login(ctx context.Context, login, password string) (LoginResult, error) {
	userID, err := s.RepoUsers.GetUserByCredentials(ctx, login, password)
	if err != nil {
		return LoginResult{}, err
	}

	_, enabled, err := s.RepoUsers.GetTOTPData(ctx, userID)
	if err != nil {
		return LoginResult{}, err
	}

	b := make([]byte, 32)
	if _, err := crypto_rand.Read(b); err != nil {
		return LoginResult{}, fmt.Errorf("serviceUsers.Login: generate token: %w", err)
	}
	token := hex.EncodeToString(b)

	// confirmed=false если 2fa включена
	if err := s.RepoUsers.SaveSession(ctx, userID, token, !enabled); err != nil {
		return LoginResult{}, err
	}

	if enabled {
		return LoginResult{TmpToken: token, Require2FA: true}, nil
	}
	fmt.Println("User authorized:", login)
	return LoginResult{Token: token, Require2FA: false}, nil
}

func (s *ServiceUsers) LoginConfirm2FA(ctx context.Context, tmpToken string, code string) (string, error) {
	userID, err := s.RepoUsers.GetUserIDByUnconfirmedToken(ctx, tmpToken)
	if err != nil {
		return "", fmt.Errorf("serviceUsers.LoginConfirm2FA: invalid or expired tmp token")
	}

	secret, _, err := s.RepoUsers.GetTOTPData(ctx, userID)
	if err != nil {
		return "", err
	}

	if !totp.Validate(code, secret) {
		return "", fmt.Errorf("serviceUsers.LoginConfirm2FA: invalid 2fa code")
	}

	if err := s.RepoUsers.ConfirmSession(ctx, tmpToken); err != nil {
		return "", err
	}

	fmt.Println("User authorized via 2fa, userID:", userID)
	return tmpToken, nil
}

func (s *ServiceUsers) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	return s.RepoUsers.GetUserIDByToken(ctx, token)
}

func (s *ServiceUsers) Setup2FA(ctx context.Context, userID int, login string) (otpauthURL string, secret string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "PepeServer",
		AccountName: login,
	})
	if err != nil {
		return "", "", fmt.Errorf("serviceUsers.Setup2FA: %w", err)
	}
	if err := s.RepoUsers.SaveTOTPSecret(ctx, userID, key.Secret()); err != nil {
		return "", "", err
	}
	return key.URL(), key.Secret(), nil
}

func (s *ServiceUsers) Confirm2FA(ctx context.Context, userID int, code string) error {
	secret, _, err := s.RepoUsers.GetTOTPData(ctx, userID)
	if err != nil {
		return err
	}
	if secret == "" {
		return fmt.Errorf("serviceUsers.Confirm2FA: secret not set, call setup first")
	}
	if !totp.Validate(code, secret) {
		return fmt.Errorf("serviceUsers.Confirm2FA: invalid code")
	}
	return s.RepoUsers.EnableTOTP(ctx, userID)
}

func (s *ServiceUsers) Delete2FA(ctx context.Context, userID int, code string) error {
	secret, enabled, err := s.RepoUsers.GetTOTPData(ctx, userID)
	if err != nil {
		return err
	}
	if !enabled {
		return fmt.Errorf("serviceUsers.Delete2FA: 2fa is not enabled")
	}
	if !totp.Validate(code, secret) {
		return fmt.Errorf("serviceUsers.Delete2FA: invalid code")
	}
	return s.RepoUsers.DisableTOTP(ctx, userID)
}

func (s *ServiceUsers) GetLoginByID(ctx context.Context, userID int) (string, error) {
	return s.RepoUsers.GetLoginByID(ctx, userID)
}

// internal/users/service.go
func (s *ServiceUsers) SaveAESKey(ctx context.Context, userID int, key []byte) error {
	return s.RepoUsers.UpdateAESKey(ctx, userID, key)
}
