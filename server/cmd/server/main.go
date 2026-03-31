package main

import (
	"lab2pepe/server/internal/core/domain"
	repositoryfiles "lab2pepe/server/internal/files/repositoryFiles"
	"lab2pepe/server/internal/files/serviceFiles"
	"lab2pepe/server/internal/files/transportFiles"
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

	http.HandleFunc("/add", HandlersFiles.HandleAddMesssage)
	http.HandleFunc("/get", HandlersFiles.HandleGetMessages)

	http.ListenAndServe(":9091", nil)
}
