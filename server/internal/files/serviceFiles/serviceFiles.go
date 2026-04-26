package serviceFiles

import (
	"context"
	"fmt"
	"io"
	"lab2pepe/server/internal/core/Aes"
	"lab2pepe/server/internal/core/domain"
	repositoryfiles "lab2pepe/server/internal/files/repositoryFiles"
)

type ServiceFiles struct {
	RepoFiles *repositoryfiles.RepoFiles
}

func NewServiceFiles(repoFiles *repositoryfiles.RepoFiles) *ServiceFiles {
	return &ServiceFiles{
		RepoFiles: repoFiles,
	}
}

func (s *ServiceFiles) AddFile(ctx context.Context, login string, encryptedBytes []byte, filename string) error {
	// получаем userID по login (уже есть в контексте)
	userID := ctx.Value(domain.UserIDKey).(int)

	// берём ключ
	aesKey, err := s.RepoFiles.GetAESKeyByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get aes key: %w", err)
	}

	// расшифровываем
	decrypted, err := Aes.Decrypt(aesKey, encryptedBytes)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	// проверяем хэш
	fileBytes, err := Aes.VerifyAndStrip(decrypted)

	if err != nil {
		return fmt.Errorf("integrity check: %w", err)
	}

	// сохраняем — передаём чистые байты и имя файла
	_, err = s.RepoFiles.SaveFile(login, fileBytes, filename)
	if err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	// сохраняем метаданные
	err = s.RepoFiles.SaveMeta(ctx, login, filename, int64(len(fileBytes)))
	if err != nil {
		return fmt.Errorf("save meta: %w", err)
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

func (s *ServiceFiles) ServiceDownload(ctx context.Context, userID int, filename string, author string) ([]byte, error) {
	// получаем файл из repo
	file, err := s.RepoFiles.RepoDownloadFile(ctx, userID, filename, author)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// читаем байты
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// берём AES-ключ
	aesKey, err := s.RepoFiles.GetAESKeyByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get aes key: %w", err)
	}

	// добавляем хэш и шифруем
	withHash := Aes.AppendHash(fileBytes)
	encrypted, err := Aes.Encrypt(aesKey, withHash)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	return encrypted, nil
}

func (s *ServiceFiles) ServiceGiveRights(ctx context.Context, user_id int, filename string, user string, rights int) error {
	if !(rights >= 1 && rights <= 3) {
		return fmt.Errorf("Error because there is not same rights")
	}
	return s.RepoFiles.RepoGiveRights(ctx, user_id, filename, user, rights)
}
