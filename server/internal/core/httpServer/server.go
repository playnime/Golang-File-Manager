package httpServer

import (
	"errors"
	"lab2pepe/server/internal/files/transportFiles"
	"lab2pepe/server/internal/users/transportUsers"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpFilesHandlers *transportFiles.HandlersFiles
	httpUsersHandlers *transportUsers.HandlersUsers

	//handlers_USER files.handlersfiles
	//handlers_FILES user.handlersuser
}

func NewHTTPServer(httpFilesHandlers *transportFiles.HandlersFiles, httpUsersHandlers *transportUsers.HandlersUsers) *HTTPServer {
	return &HTTPServer{
		httpFilesHandlers: httpFilesHandlers,
		httpUsersHandlers: httpUsersHandlers,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.HandleFunc("/add", s.httpFilesHandlers.HandleAddMesssage).Methods("POST")
	router.HandleFunc("/get", s.httpFilesHandlers.HandleGetMessages).Methods("GET")
	router.HandleFunc("/addusr", s.httpUsersHandlers.HandleAddUser).Methods("POST")

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
