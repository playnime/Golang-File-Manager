package main

import (
	"lab2pepe/server/internal/core/domain"
	repositoryfiles "lab2pepe/server/internal/files/repositoryFiles"
	"lab2pepe/server/internal/files/serviceFiles"
	"lab2pepe/server/internal/files/transportFiles"
	"lab2pepe/server/internal/users/repositoryUsers"
	"lab2pepe/server/internal/users/serviceUsers"
	"lab2pepe/server/internal/users/transportUsers"
	"net/http"
)

func main() {
	Repo := domain.NewRepo()
	Ilya := domain.User{
		Login:    "Ilya",
		Password: "123",
		Messages: nil,
	}
	Repo.Bd = append(Repo.Bd, Ilya)

	RepoFiles := repositoryfiles.NewHRepoFiles(Repo)
	ServiceFiles := serviceFiles.NewServiceFiles(RepoFiles)
	HandlersFiles := transportFiles.NewHandlersFiles(ServiceFiles)

	RepoUsers := repositoryUsers.NewRepoUsers(Repo)
	ServiceUsers := serviceUsers.NewServiceUsers(RepoUsers)
	HandlersUsers := transportUsers.NewHandlersUser(ServiceUsers)

	http.HandleFunc("/add", HandlersFiles.HandleAddMesssage)
	http.HandleFunc("/get", HandlersFiles.HandleGetMessages)
	http.HandleFunc("/addusr", HandlersUsers.HandleAddUser)

	http.ListenAndServe(":9091", nil)
}
