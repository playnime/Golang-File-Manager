package serviceFiles

import (
	"fmt"
	repositoryfiles "lab2pepe/server/internal/files/repositoryFiles"
	"mime/multipart"
)

type ServiceFiles struct {
	RepoFiles *repositoryfiles.RepoFiles
}

func NewServiceFiles(repoFiles *repositoryfiles.RepoFiles) *ServiceFiles {
	return &ServiceFiles{
		RepoFiles: repoFiles,
	}
}

/*func (s *ServiceFiles) ServiceAddFile(login string, password string, message string) error {
	return s.RepoFiles.RepoAddFile(login, password, message)
}*/

func (s *ServiceFiles) AddFile(login string, password string, file multipart.File, header *multipart.FileHeader) error {
	path, err := s.RepoFiles.SaveFile(login, file, header.Filename)
	if err != nil {
		fmt.Println("error on SaveFile: ", err)
		return err
	}
	_ = path
	err = s.RepoFiles.SaveMeta(login, password, header)
	if err != nil {
		fmt.Println("error on SaveMeta: ", err)
		return err
	}
	return nil
}

func (s *ServiceFiles) ServiceGetFiles(login string, password string) ([]string, error) {
	return s.RepoFiles.RepoGetFiles(login, password)
}
