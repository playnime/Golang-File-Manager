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

func NewRepo() *Repository {
	return &Repository{
		Bd: make([]User, 0),
	}
}
