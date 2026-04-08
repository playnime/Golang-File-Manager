package transportFiles

import (
	"encoding/json"
	"errors"
	"fmt"
	"lab2pepe/server/internal/core/domain"
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

func (h *HandlersFiles) HandleAddFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println("error on parsing paket: ", err)
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("File")
	if err != nil {
		fmt.Println("error on parsing file: ", err)
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	id := r.Context().Value(domain.UserIDKey).(int)
	login, err := h.ServiceFiles.GetLoginById(r.Context(), int(id))
	if err != nil {
		fmt.Println("handler: error on founding login: ", err)
		httpUtils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Println("login: ", login)
	fmt.Println("Filename: ", header.Filename)
	fmt.Println("Filesize: ", header.Size)
	fmt.Println("mime: ", header.Header.Get("Content-Type"))

	err = h.ServiceFiles.AddFile(r.Context(), login, file, header)
	if err != nil {
		fmt.Println("error on creating qyery in database")
		httpUtils.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("file is added"))
}

func (h *HandlersFiles) HandleGetFiles(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(domain.UserIDKey).(int)
	if !ok {
		httpUtils.WriteError(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	login, err := h.ServiceFiles.GetLoginById(r.Context(), id)
	if err != nil {
		fmt.Println("error on getting login by id:", err)
		httpUtils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	files, err := h.ServiceFiles.ServiceGetFiles(r.Context(), login)
	if err != nil {
		fmt.Println("handler: error:", err)
		httpUtils.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(GetFilesResponseDTO{Files: files}); err != nil {
		fmt.Println("handler: encode error:", err)
	}
}

/*func (h *HandlersFiles) HandleGetFiles(w http.ResponseWriter, r *http.Request) {
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
*/
