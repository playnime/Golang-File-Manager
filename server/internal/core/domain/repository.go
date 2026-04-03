package domain

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Repository получает *sql.DB (соединение с БД)

type User struct {
	Login    string
	Password string
	Messages []string
}

type Repository struct {
	Bd   []User
	Conn *pgx.Conn
	Ctx  context.Context
}

func NewRepo(ctx context.Context, conn *pgx.Conn) *Repository {
	return &Repository{
		Bd:   make([]User, 0),
		Ctx:  ctx,
		Conn: conn,
	}
}
