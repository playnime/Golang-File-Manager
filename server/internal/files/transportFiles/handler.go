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

func (h *HandlersFiles) HandleAddFile(w http.ResponseWriter, r *http.Request) {
	// h.AddFile
}
