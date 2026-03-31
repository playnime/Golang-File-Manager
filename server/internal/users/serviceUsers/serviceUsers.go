package serviceUsers

import (
	"lab2pepe/server/internal/users/repositoryUsers"
)

type ServiceFiles struct {
	RepoFiles *repositoryUsers.RepoUsers
}

func NewServiceUsers() {

}

func (s *ServiceFiles) AddUsers() {
	// s.CreateFiles
}
