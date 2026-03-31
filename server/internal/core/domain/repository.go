package domain

// Repository получает *sql.DB (соединение с БД)

type User struct {
	Login    string
	Password string
	Messages []string
}

type Repository struct {
	Bd []User
}
