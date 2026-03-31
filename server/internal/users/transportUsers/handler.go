package transportUsers

import (
	"encoding/json"
	"fmt"
	"lab2pepe/server/internal/users/serviceUsers"
	"net/http"
)

type HandlersUsers struct {
	ServiceFiles *serviceUsers.ServiceFiles
}

func NewHandlersUser() {

}

func (h *HandlersUsers) HandleAddUser(w http.ResponseWriter, r *http.Request) {

	var user AddUsersDTO
	if err := json.NewDecoder(r.Body).Decode(&AddUsersDTO); err != nil {
		fmt.Println("err:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
