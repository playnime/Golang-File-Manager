package transportFiles

type AddMessageDTO struct {
	Login    string `json:"Login"`
	Password string `json:"Password"`
	Message  string `json:"Message"`
}

type GetMessagesDTO struct {
	Login    string `json:"Login"`
	Password string `json:"Password"`
}

type GetMessagesResponseDTO struct {
	Messages []string
}
