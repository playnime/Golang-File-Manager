package transportFiles

import (
	"encoding/json"
	"fmt"
	"io"
	"lab2pepe/server/internal/core/httpUtils"
	"lab2pepe/server/internal/files/serviceFiles"
	"net/http"
)

type HandlersFiles struct {
	ServiceFiles *serviceFiles.ServiceFiles
}

func NewHandlersFiles(serviceFiles *serviceFiles.ServiceFiles) *HandlersFiles {
	return &HandlersFiles{
		ServiceFiles: serviceFiles,
	}
}

func (h *HandlersFiles) HandleAddMesssage(w http.ResponseWriter, r *http.Request) {
	httpResonce, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error on reading body: ", err)
		httpUtils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	var User AddFileDTO

	err = json.Unmarshal(httpResonce, &User)
	if err != nil {
		fmt.Println("error on unmarshalling response body: ", err)
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err = h.ServiceFiles.ServiceAddFile(User.Login, User.Password, User.File)
	if err != nil {
		fmt.Println("error on adding file: ", err)
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}
	fmt.Println("file is on database")
	w.WriteHeader(http.StatusCreated)
}

func (h *HandlersFiles) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	httpResonce, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error on reading body: ", err)
		//w.WriteHeader(http.StatusInternalServerError)
		httpUtils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	var User GetFilesDTO

	err = json.Unmarshal(httpResonce, &User)
	if err != nil {
		fmt.Println("error on unmarshalling response body: ", err)
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	files, err := h.ServiceFiles.ServiceGetFiles(User.Login, User.Password)
	if err != nil {
		fmt.Println("error on getting files: ", err)
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	response := GetFilesResponseDTO{
		Files: files,
	}

	b, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error on marshaling http answer: ", err)
		httpUtils.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		fmt.Println("error on write http answer: ", err)
		return
	}

	fmt.Println("User gets files")
}
