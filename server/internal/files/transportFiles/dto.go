package transportFiles

import "lab2pepe/server/internal/core/domain"

type AddFileDTO struct {
	Login    string `json:"Login"`
	Password string `json:"Password"`
	File     string `json:"File"`
}

type GetFilesDTO struct {
	Login    string `json:"Login"`
	Password string `json:"Password"`
}

type GetFilesResponseDTO struct {
	Files []domain.GetFilesModel `json:"Files"`
}
