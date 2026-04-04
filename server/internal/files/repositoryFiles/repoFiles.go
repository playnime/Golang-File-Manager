package repositoryfiles

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoFiles struct {
	pool *pgxpool.Pool
}

func NewHRepoFiles(pool *pgxpool.Pool) *RepoFiles {
	return &RepoFiles{
		pool: pool,
	}
}

func (r *RepoFiles) SaveFile(login string, file multipart.File, filename string) (string, error) {
	dir := filepath.Join(os.Getenv("UPLOAD_DIR"), login)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println("error on creating directory: ", err)
		return "", err
	}

	path := filepath.Join(dir, filename)
	data, err := os.Create(path)
	if err != nil {
		fmt.Println("error on creating file: ", err)
		return "", err
	}
	defer data.Close()

	_, err = io.Copy(data, file)
	if err != nil {
		fmt.Println("error on copying file: ", err)
		return "", err
	}

	return path, nil
}

func (r *RepoFiles) SaveMeta(login string, password string, header *multipart.FileHeader) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	path := filepath.Join(login, header.Filename)

	_, err := r.pool.Exec(ctx,
		`INSERT INTO files (user_id, file_name, path, size_bytes, mime_type, gd)
         SELECT id, $3, $4, $5, $6, $7 FROM users
         WHERE login = $1 AND password = $2`,
		login, password, header.Filename, path, header.Size, header.Header.Get("Content-Type"), 3,
	)
	if err != nil {
		fmt.Println("Error on query: ", err)
		return err
	}

	fmt.Println("lucky!!")
	return nil
}

func (r *RepoFiles) RepoGetFiles(login string, password string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx,
		`SELECT f.file_name FROM files f
         JOIN users u ON u.id = f.user_id
         WHERE u.login = $1 AND u.password = $2`,
		login, password,
	)
	if err != nil {
		return nil, fmt.Errorf("repoFiles.RepoGetFiles: %w", err)
	}
	defer rows.Close()

	var files []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("repoFiles.RepoGetFiles scan: %w", err)
		}
		files = append(files, name)
	}

	return files, nil
}
