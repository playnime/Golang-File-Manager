package serviceUsers

import (
	"lab2pepe/server/internal/users/repositoryUsers"
)

type ServiceUsers struct {
	RepoFiles *repositoryUsers.RepoUsers
}

func NewServiceUsers(service *ServiceUsers) {

}

func (s *ServiceUsers) AddUsers() {
	// s.CreateUsers
}
