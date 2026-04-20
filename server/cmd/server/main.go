package main

import (
	"context"
	"fmt"
	"lab2pepe/server/internal/core/db"
	"lab2pepe/server/internal/core/domain"
	"lab2pepe/server/internal/core/httpServer"
	repositoryfiles "lab2pepe/server/internal/files/repositoryFiles"
	"lab2pepe/server/internal/files/serviceFiles"
	"lab2pepe/server/internal/files/transportFiles"
	"lab2pepe/server/internal/users/repositoryUsers"
	"lab2pepe/server/internal/users/serviceUsers"
	"lab2pepe/server/internal/users/transportUsers"
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
	if err := db.AutoMigrate(ctx, Pool); err != nil {
		fmt.Println("migration error:", err)
		os.Exit(1)
	}

	RepoFiles := repositoryfiles.NewHRepoFiles(Pool)
	ServiceFiles := serviceFiles.NewServiceFiles(RepoFiles)
	HandlersFiles := transportFiles.NewHandlersFiles(ServiceFiles)

	RepoUsers := repositoryUsers.NewRepoUsers(Pool)
	ServiceUsers := serviceUsers.NewServiceUsers(RepoUsers)
	HandlersUsers := transportUsers.NewHandlersUser(ServiceUsers)

	server := httpServer.NewHTTPServer(HandlersFiles, HandlersUsers, ServiceUsers)
	if err := server.StartServer(); err != nil {
		panic(err)
	}
}
