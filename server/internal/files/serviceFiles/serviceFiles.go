package serviceFiles

import repositoryfiles "lab2pepe/server/internal/files/repositoryFiles"

type ServiceFiles struct {
	RepoFiles *repositoryfiles.RepoFiles
}

func NewServiceFiles() {

}

func (s *ServiceFiles) AddFile() {
	// s.CreateFiles
}
