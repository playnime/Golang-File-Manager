package serviceFiles

import repositoryfiles "lab2pepe/server/internal/files/repositoryFiles"

type ServiceFiles struct {
	RepoFiles *repositoryfiles.RepoFiles
}

func NewServiceFiles(repoFiles *repositoryfiles.RepoFiles) *ServiceFiles {
	return &ServiceFiles{
		RepoFiles: repoFiles,
	}
}

func (s *ServiceFiles) AddMessage(login string, password string, message string) {
	s.RepoFiles.CreateMessage(login, password, message)
}

func (s *ServiceFiles) GetMessages(login string, password string) []string {
	return s.RepoFiles.GetMessages(login, password)
}
