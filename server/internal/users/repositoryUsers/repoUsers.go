package repositoryUsers

import (
	"context"
	"fmt"
	"lab2pepe/server/internal/core/md5MY"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoUsers struct {
	Pool *pgxpool.Pool
}

func NewRepoUsers(pool *pgxpool.Pool) *RepoUsers {
	return &RepoUsers{Pool: pool}
}

func (r *RepoUsers) CreateUser(ctx context.Context, login, password string) error {
	hash := md5MY.HashToString(md5MY.HASH([]byte(password)))
	_, err := r.Pool.Exec(ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2)`,
		login, hash,
	)
	if err != nil {
		return fmt.Errorf("repoUsers.CreateUser: %w", err)
	}
	fmt.Println("User created:", login)
	return nil
}

func (r *RepoUsers) GetUserByCredentials(ctx context.Context, login, password string) (int, error) {
	hash := md5MY.HashToString(md5MY.HASH([]byte(password)))
	var userID int
	err := r.Pool.QueryRow(ctx,
		`SELECT id FROM users WHERE login = $1 AND password = $2`,
		login, hash,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("repoUsers.GetUserByCredentials: %w", err)
	}
	return userID, nil
}

func (r *RepoUsers) SaveSession(ctx context.Context, userID int, token string, confirmed bool) error {
	_, err := r.Pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at < NOW()`)
	if err != nil {
		return fmt.Errorf("repoUsers.SaveSession: clear: %w", err)
	}
	_, err = r.Pool.Exec(ctx,
		`INSERT INTO sessions (user_id, token, expires_at, confirmed)
		 VALUES ($1, $2, NOW() + INTERVAL '5 minutes', $3)`,
		userID, token, confirmed,
	)
	if err != nil {
		return fmt.Errorf("repoUsers.SaveSession: %w", err)
	}
	return nil
}

func (r *RepoUsers) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	var userID int
	err := r.Pool.QueryRow(ctx,
		`SELECT user_id FROM sessions
		 WHERE token = $1 AND expires_at > NOW() AND confirmed = TRUE`,
		token,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("repoUsers.GetUserIDByToken: %w", err)
	}
	return userID, nil
}

func (r *RepoUsers) GetUserIDByUnconfirmedToken(ctx context.Context, token string) (int, error) {
	var userID int
	err := r.Pool.QueryRow(ctx,
		`SELECT user_id FROM sessions
		 WHERE token = $1 AND expires_at > NOW() AND confirmed = FALSE`,
		token,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("repoUsers.GetUserIDByUnconfirmedToken: %w", err)
	}
	return userID, nil
}

func (r *RepoUsers) ConfirmSession(ctx context.Context, token string) error {
	tag, err := r.Pool.Exec(ctx,
		`UPDATE sessions SET confirmed = TRUE
		 WHERE token = $1 AND expires_at > NOW()`,
		token,
	)
	if err != nil {
		return fmt.Errorf("repoUsers.ConfirmSession: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("repoUsers.ConfirmSession: token not found or expired")
	}
	return nil
}

func (r *RepoUsers) GetTOTPData(ctx context.Context, userID int) (secret string, enabled bool, err error) {
	err = r.Pool.QueryRow(ctx,
		`SELECT totp_secret, totp_enabled FROM users WHERE id = $1`,
		userID,
	).Scan(&secret, &enabled)
	if err != nil {
		return "", false, fmt.Errorf("repoUsers.GetTOTPData: %w", err)
	}
	return secret, enabled, nil
}

func (r *RepoUsers) SaveTOTPSecret(ctx context.Context, userID int, secret string) error {
	_, err := r.Pool.Exec(ctx,
		`UPDATE users SET totp_secret = $1, totp_enabled = FALSE WHERE id = $2`,
		secret, userID,
	)
	if err != nil {
		return fmt.Errorf("repoUsers.SaveTOTPSecret: %w", err)
	}
	return nil
}

func (r *RepoUsers) EnableTOTP(ctx context.Context, userID int) error {
	_, err := r.Pool.Exec(ctx,
		`UPDATE users SET totp_enabled = TRUE WHERE id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("repoUsers.EnableTOTP: %w", err)
	}
	return nil
}

func (r *RepoUsers) DisableTOTP(ctx context.Context, userID int) error {
	_, err := r.Pool.Exec(ctx,
		`UPDATE users SET totp_secret = '', totp_enabled = FALSE WHERE id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("repoUsers.DisableTOTP: %w", err)
	}
	return nil
}

func (r *RepoUsers) GetLoginByID(ctx context.Context, userID int) (string, error) {
	var login string
	err := r.Pool.QueryRow(ctx,
		`SELECT login FROM users WHERE id = $1`,
		userID,
	).Scan(&login)
	if err != nil {
		return "", fmt.Errorf("repoUsers.GetLoginByID: %w", err)
	}
	return login, nil
}
