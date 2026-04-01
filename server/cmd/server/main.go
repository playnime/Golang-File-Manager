package main

import (
	"lab2pepe/server/internal/core/domain"
	"lab2pepe/server/internal/core/httpServer"
	repositoryfiles "lab2pepe/server/internal/files/repositoryFiles"
	"lab2pepe/server/internal/files/serviceFiles"
	"lab2pepe/server/internal/files/transportFiles"
	"lab2pepe/server/internal/users/repositoryUsers"
	"lab2pepe/server/internal/users/serviceUsers"
	"lab2pepe/server/internal/users/transportUsers"
	"log"
)

func main() {
	Repo := domain.NewRepo()
	Repo.Bd = append(Repo.Bd, domain.User{
		Login:    "Ilya",
		Password: "123",
	})

	RepoFiles := repositoryfiles.NewHRepoFiles(Repo)
	ServiceFiles := serviceFiles.NewServiceFiles(RepoFiles)
	HandlersFiles := transportFiles.NewHandlersFiles(ServiceFiles)

	RepoUsers := repositoryUsers.NewRepoUsers(Repo)
	ServiceUsers := serviceUsers.NewServiceUsers(RepoUsers)
	HandlersUsers := transportUsers.NewHandlersUser(ServiceUsers)

	server := httpServer.NewHTTPServer(HandlersFiles, HandlersUsers)
	if err := server.StartServer(); err != nil {
		log.Fatal(err)
	}
}
