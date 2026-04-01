package transportUsers

import "errors"

type AddUsersDTO struct {
	Login string `json:"login"`

	Password string `json:"password"`
}

func (t AddUsersDTO) ValidateForCreateUser() error {
	if len(t.Login) < 4 || len(t.Login) > 20 {
		return errors.New("Login not correct")
	}

	if len(t.Password) < 4 || len(t.Password) > 20 {
		return errors.New("password is not correct")
	}

	return nil
}

type CreateUserResponseDTO struct {
	Login string `json:"login"`
}
