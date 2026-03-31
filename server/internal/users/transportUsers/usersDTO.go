package transportUsers

type AddUsersDTO struct {
	//Описание покупки
	Login string `json:"login"`

	//Сумма покупки
	Password int `json:"password"`
}
