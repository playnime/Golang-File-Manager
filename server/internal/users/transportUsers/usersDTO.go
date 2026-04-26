package transportUsers

import "errors"

type AddUsersDTO struct {
	Login    string `json:"login"`
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

type LoginDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Token      string `json:"token,omitempty"`
	TmpToken   string `json:"tmp_token,omitempty"`
	Require2FA bool   `json:"require_2fa"`
}

type Login2FADTO struct {
	TmpToken string `json:"tmp_token"`
	Code     string `json:"code"`
}

type Setup2FAResponseDTO struct {
	OtpauthURL string `json:"otpauth_url"`
	Secret     string `json:"secret"`
}

type Confirm2FADTO struct {
	Code string `json:"code"`
}

type Delete2FADTO struct {
	Code string `json:"code"`
}

type keyExchangeRequest struct {
	N string `json:"n"`
	E string `json:"e"`
}

type keyExchangeResponse struct {
	EncryptedAESKey string `json:"encrypted_aes_key"`
}
