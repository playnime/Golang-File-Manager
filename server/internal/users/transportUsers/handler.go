package transportUsers

import (
	"encoding/json"
	"fmt"
	"lab2pepe/server/internal/core/httpUtils"
	"lab2pepe/server/internal/users/serviceUsers"
	"net/http"
)

type HandlersUsers struct {
	ServiceFiles *serviceUsers.ServiceUsers
}

func NewHandlersUser(service *serviceUsers.ServiceUsers) *HandlersUsers {
	return &HandlersUsers{ServiceFiles: service}
}

func (h *HandlersUsers) HandleAddUser(w http.ResponseWriter, r *http.Request) {

	var dto AddUsersDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := dto.ValidateForCreateUser(); err != nil {
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := h.ServiceFiles.AddUsers(dto.Login, dto.Password); err != nil {
		httpUtils.WriteError(w, err, http.StatusConflict)
		return
	}

	response := CreateUserResponseDTO{
		Login: dto.Login,
	}

	b, err := json.MarshalIndent(response, "", "    ")

	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}

}
