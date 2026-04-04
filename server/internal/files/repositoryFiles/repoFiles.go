package repositoryfiles

import (
	"context"
	"fmt"
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

func (r *RepoFiles) RepoAddFile(login string, password string, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	path := login + "/" + message

	_, err := r.pool.Exec(ctx,
		`INSERT INTO files (user_id, file_name, path)
         SELECT id, $3, $4 FROM users
         WHERE login = $1 AND password = $2`,
		login, password, message, path,
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

// func (r *RepoFiles) RepoAddFile(login string, password string, message string) error {
// 	var k int = -1
// 	for i, v := range r.repo.Bd {
// 		if login == v.Login {
// 			if password == v.Password {
// 				k = i
// 			} else {
// 				return errors.New("password is not correct")
// 			}
// 			break
// 		}
// 	}

// 	if k != -1 {
// 		r.repo.Bd[k].Messages = append(r.repo.Bd[k].Messages, message)
// 	} else {
// 		return errors.New("Login is not correct")
// 	}

// 	return nil
// }

// func (r *RepoFiles) RepoGetFiles(login string, password string) ([]string, error) {
// 	var k int = -1
// 	for i, v := range r.repo.Bd {
// 		if login == v.Login {
// 			if password == v.Password {
// 				k = i
// 			} else {
// 				return nil, errors.New("password is not correct")
// 			}
// 			break
// 		}
// 	}

// 	if k != -1 {
// 		return r.repo.Bd[k].Messages, nil
// 	}

// 	return nil, errors.New("Login is not correct")
// }

//еще будет функция скачивания файла
