package transportFiles

import (
	"encoding/json"
	"fmt"
	"io"
	"lab2pepe/server/internal/files/serviceFiles"
	"net/http"
)

type HandlersFiles struct {
	ServiceFiles *serviceFiles.ServiceFiles
}

func NewHandlersFiles(serviceFiles *serviceFiles.ServiceFiles) *HandlersFiles {
	return &HandlersFiles{
		ServiceFiles: serviceFiles,
	}
}

func (h *HandlersFiles) HandleAddMesssage(w http.ResponseWriter, r *http.Request) {
	httpResonce, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error on reading body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var User AddMessageDTO

	err = json.Unmarshal(httpResonce, &User)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("error on unmarshalling response body: ", err)
		return
	}

	h.ServiceFiles.AddMessage(User.Login, User.Password, User.Message)
	fmt.Println("Message is on slice")
}

func (h *HandlersFiles) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	httpResonce, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error on reading body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var User GetMessagesDTO

	err = json.Unmarshal(httpResonce, &User)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("error on unmarshalling response body: ", err)
		return
	}

	response := GetMessagesResponseDTO{
		Messages: h.ServiceFiles.GetMessages(User.Login, User.Password),
	}

	b, err := json.Marshal(response)
	if err != nil {
		fmt.Println("не удалось зашифровать строку: ", err)
	}
	w.Write(b)
}
