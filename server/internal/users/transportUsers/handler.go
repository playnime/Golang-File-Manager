package transportUsers

import (
	"encoding/json"
	"lab2pepe/server/internal/core/httpServer"
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
		httpServer.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := h.ServiceFiles.AddUsers(dto.Login, dto.Password); err != nil {
		httpServer.WriteError(w, err, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
