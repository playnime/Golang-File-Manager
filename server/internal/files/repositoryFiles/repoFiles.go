package repositoryfiles

import (
	"context"
	"errors"
	"fmt"
	"io"
	"lab2pepe/server/internal/core/domain"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
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

func (r *RepoFiles) SaveMeta(ctx context.Context, login string, header *multipart.FileHeader) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	path := filepath.Join(login, header.Filename)

	_, err := r.pool.Exec(ctx,
		`INSERT INTO files (user_id, file_name, path, size_bytes, mime_type, gd)
         SELECT id, $2, $3, $4, $5, $6 FROM users
         WHERE login = $1`,
		login, header.Filename, path, header.Size, header.Header.Get("Content-Type"), 3,
	)
	if err != nil {
		fmt.Println("Error on query: ", err)
		return err
	}

	fmt.Println("lucky!!")
	return nil
}

func (r *RepoFiles) RepoGetFiles(ctx context.Context, id int) ([]domain.GetFilesModel, error) {
	files := make([]domain.GetFilesModel, 0)
	rows, err := r.pool.Query(ctx,
		`SELECT file_name, size_bytes, gd, uploaded_at
		FROM files
		WHERE user_id = $1`,
		id,
	)

	if err != nil {
		fmt.Println("repo: err on getting rows: ", err)
		return nil, err
	}

	defer rows.Close()

	var file domain.GetFilesModel

	for rows.Next() {
		err := rows.Scan(&file.FileName, &file.SizeBytes, &file.Gd, &file.UploadedAt)
		if err != nil {
			fmt.Println("repo: error on scanning rows: ", err)
			return nil, err
		}

		files = append(files, file)
	}

	fmt.Println("client getting files")
	return files, nil
}

func (r *RepoFiles) GetLoginById(ctx context.Context, id int) (string, error) {
	var login string
	err := r.pool.QueryRow(ctx,
		`SELECT login FROM users WHERE id = $1`,
		id,
	).Scan(&login)
	if err != nil {
		fmt.Println("didn1t found user with that id: ", err)
		return "", err
	}
	fmt.Println("repo: login: ", login)

	return login, nil
}

func (r *RepoFiles) RepoDeleteFiles(ctx context.Context, user_id int, filename string, author string) error {
	var f DeleteModel

	err := r.pool.QueryRow(ctx,
		`SELECT id, file_name, path, gd
		FROM files
		WHERE user_id = $1 AND file_name = $2 AND path = $3`,
		user_id, filename, author+"/"+filename,
	).Scan(&f.id, &f.filename, &f.filepath, &f.gd)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("file not found")
		} else {
			fmt.Println("error on selecting fileMeta: ", err)
		}
		return err
	}

	if f.gd != 3 && f.gd != 1 {
		fmt.Println("You dont have permission")
		return errors.New("You dont have permission")
	}

	_, err = r.pool.Exec(
		ctx,
		`DELETE FROM files
		WHERE path = $1`,
		f.filepath,
	)
	if err != nil {
		fmt.Println("error on deleting meta: ", err)
		return err
	}

	if !filepath.IsLocal(f.filepath) {
		return fmt.Errorf("invalid file path")
	}

	err = os.Remove("out/" + f.filepath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("error on deleting file on disk: ", err)
		return err
	}

	fmt.Println("deleting is correct!")
	return nil
}

func (r *RepoFiles) RepoDownloadFile(ctx context.Context, user_id int, filename string, author string) (io.ReadCloser, error) {
	var f DownloadModel

	err := r.pool.QueryRow(ctx,
		`SELECT id, file_name, path, gd
		FROM files
		WHERE user_id = $1 AND file_name = $2 AND path = $3`,
		user_id, filename, author+"/"+filename,
	).Scan(&f.id, &f.filename, &f.filepath, &f.gd)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("file not found")
		} else {
			fmt.Println("error on selecting fileMeta: ", err)
		}
		return nil, err
	}

	if f.gd != 3 && f.gd != 2 {
		fmt.Println("You dont have permission")
		return nil, errors.New("You dont have permission")
	}

	file, err := os.Open("out/" + f.filepath)
	if err != nil {
		fmt.Println("error on opening file: ", err)
		return nil, err
	}

	return file, nil
}
