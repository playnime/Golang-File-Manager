package domain

// Repository получает *sql.DB (соединение с БД)

type User struct {
	login    string
	password int
	messages []string
}

type Repository struct {
	bd []User
}
