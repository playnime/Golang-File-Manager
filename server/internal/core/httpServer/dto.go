package httpServer

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func WriteError(w http.ResponseWriter, err error, status int) {
	errDTO := ErrorDTO{Message: err.Error(), Time: time.Now()}
	http.Error(w, errDTO.ToString(), status)
}
