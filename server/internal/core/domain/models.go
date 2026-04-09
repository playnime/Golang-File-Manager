package domain

import "time"

type GetFilesModel struct {
	FileName   string    `json:"file_name"`
	SizeBytes  int64     `json:"size_bytes"`
	Gd         int       `json:"gd"`
	UploadedAt time.Time `json:"uploaded_at"`
}
