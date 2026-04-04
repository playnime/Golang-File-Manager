package main

import (
	"context"
	"lab2pepe/server/internal/core/domain"
	"lab2pepe/server/internal/core/httpServer"
	repositoryfiles "lab2pepe/server/internal/files/repositoryFiles"
	"lab2pepe/server/internal/files/serviceFiles"
	"lab2pepe/server/internal/files/transportFiles"
	"os"
)

func main() {
	// ctx := context.Background()

	// conn_str := os.Getenv("CONN_STR")
	// conn, err := pgx.Connect(ctx, conn_str)
	// if err != nil {
	// 	panic(err)
	// }

	ctx := context.Background()

	Pool, err := domain.NewPool(ctx, os.Getenv("CONN_STR"))
	if err != nil {
		panic(err)
	}
	defer Pool.Close()

	Repo := domain.NewRepo()
	Repo.Bd = append(Repo.Bd, domain.User{
		Login:    "Ilya",
		Password: "123",
	})

	RepoFiles := repositoryfiles.NewHRepoFiles(Pool)
	ServiceFiles := serviceFiles.NewServiceFiles(RepoFiles)
	HandlersFiles := transportFiles.NewHandlersFiles(ServiceFiles)

	// RepoUsers := repositoryUsers.NewRepoUsers(Repo)
	// ServiceUsers := serviceUsers.NewServiceUsers(RepoUsers)
	// HandlersUsers := transportUsers.NewHandlersUser(ServiceUsers)

	server := httpServer.NewHTTPServer(HandlersFiles, nil)
	if err := server.StartServer(); err != nil {
		panic(err)
	}
}
