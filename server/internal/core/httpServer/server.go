package httpServer

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers

	//handlers_USER files.handlersfiles
	//handlers_FILES user.handlersuser
}

func NewHTTPServer(HTTPHandler *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: HTTPHandler,
	}
}

func (s *HTTPServer) StartServer() error {

	router := mux.NewRouter()

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
