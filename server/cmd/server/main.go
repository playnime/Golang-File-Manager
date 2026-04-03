package main

import (
	"context"
	"lab2pepe/server/internal/core/domain"
	"lab2pepe/server/internal/core/httpServer"
	repositoryfiles "lab2pepe/server/internal/files/repositoryFiles"
	"lab2pepe/server/internal/files/serviceFiles"
	"lab2pepe/server/internal/files/transportFiles"
	"lab2pepe/server/internal/users/repositoryUsers"
	"lab2pepe/server/internal/users/serviceUsers"
	"lab2pepe/server/internal/users/transportUsers"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	conn_str := os.Getenv("CONN_STR")
	conn, err := pgx.Connect(ctx, conn_str)
	if err != nil {
		panic(err)
	}

	Repo := domain.NewRepo(ctx, conn)
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
		panic(err)
	}
}
