package transportFiles

import (
	"lab2pepe/server/internal/files/serviceFiles"
	"net/http"
)

type HandlersFiles struct {
	ServiceFiles *serviceFiles.ServiceFiles
}

func NewHandlersFiles() {

}

func (h *HandlersFiles) HandlersFiles(w http.ResponseWriter, r *http.Request) {

}
