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

func (s *ServiceFiles) ServiceAddFile(login string, password string, message string) error {
	return s.RepoFiles.RepoAddFile(login, password, message)
}

func (s *ServiceFiles) ServiceGetFiles(login string, password string) ([]string, error) {
	return s.RepoFiles.RepoGetFiles(login, password)
}
