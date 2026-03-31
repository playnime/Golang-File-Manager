package transportUsers

import (
	"lab2pepe/server/internal/users/serviceUsers"
	"net/http"
)

type HandlersUsers struct {
	ServiceFiles *serviceUsers.ServiceFiles
}

func NewHandlersUser() {

}

func (h *HandlersUsers) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	// h.AddFile
}
