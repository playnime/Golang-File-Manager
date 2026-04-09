package httpServer

import (
	"context"
	"errors"
	"fmt"
	"lab2pepe/server/internal/core/domain"
	"lab2pepe/server/internal/files/transportFiles"
	"lab2pepe/server/internal/users/transportUsers"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type AuthService interface {
	GetUserIDByToken(ctx context.Context, token string) (int, error)
}

type HTTPServer struct {
	httpFilesHandlers *transportFiles.HandlersFiles
	httpUsersHandlers *transportUsers.HandlersUsers
	authService       AuthService
}

func NewHTTPServer(
	httpFilesHandlers *transportFiles.HandlersFiles,
	httpUsersHandlers *transportUsers.HandlersUsers,
	authService AuthService,
) *HTTPServer {
	return &HTTPServer{
		httpFilesHandlers: httpFilesHandlers,
		httpUsersHandlers: httpUsersHandlers,
		authService:       authService,
	}
}

func (s *HTTPServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			fmt.Println("unauthorized, no prefix")
			http.Error(w, "unauthorized, no prefix", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")

		userID, err := s.authService.GetUserIDByToken(r.Context(), token)
		if err != nil {
			fmt.Println("unauthorized: ", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), domain.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.HandleFunc("/register", s.httpUsersHandlers.HandleAddUser).Methods("POST")
	router.HandleFunc("/login", s.httpUsersHandlers.HandleLogin).Methods("POST")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(s.authMiddleware)
	protected.HandleFunc("/add", s.httpFilesHandlers.HandleAddFile).Methods("POST")
	protected.HandleFunc("/get", s.httpFilesHandlers.HandleGetFiles).Methods("GET")
	protected.HandleFunc("/delete", s.httpFilesHandlers.HandleDeleteFiles).Methods("DELETE")

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
