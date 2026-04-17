package serviceFiles

import (
	"context"
	"fmt"
	"io"
	"lab2pepe/server/internal/core/domain"
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

func (s *ServiceFiles) AddFile(ctx context.Context, login string, file multipart.File, header *multipart.FileHeader) error {
	path, err := s.RepoFiles.SaveFile(login, file, header.Filename)
	if err != nil {
		fmt.Println("error on SaveFile: ", err)
		return err
	}
	_ = path
	err = s.RepoFiles.SaveMeta(ctx, login, header)
	if err != nil {
		fmt.Println("error on SaveMeta: ", err)
		return err
	}
	return nil
}

func (s *ServiceFiles) ServiceGetFiles(ctx context.Context, id int) ([]domain.GetFilesModel, error) {
	return s.RepoFiles.RepoGetFiles(ctx, id)
}

func (s *ServiceFiles) GetLoginById(ctx context.Context, id int) (string, error) {
	return s.RepoFiles.GetLoginById(ctx, id)
}

func (s *ServiceFiles) ServiceDelete(ctx context.Context, id int, filename string, author string) error {
	return s.RepoFiles.RepoDeleteFiles(ctx, id, filename, author)
}

func (s *ServiceFiles) ServiceDownload(ctx context.Context, user_id int, filename string, author string) (io.ReadCloser, error) {
	return s.RepoFiles.RepoDownloadFile(ctx, user_id, filename, author)
}
