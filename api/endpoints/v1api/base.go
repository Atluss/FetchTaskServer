package v1api

import (
	"encoding/json"
	"github.com/Atluss/FetchTaskServer/lib"
	"net/http"
)

const (
	V1ApiQueue = "v1"
)

type HeadRequest interface {
	Request()   // execute FetchTask
	NatsQueue() // nats func
}

type ApiRun interface {
	Execute()  // запуск исполняющей функции в запросе
	Validate() // валидация данных
}

type ApiRequest struct {
	HeadRequest
	w *http.ResponseWriter
	r *http.Request
}

// Bad request
type ReplayBadRequest struct {
	Status      int    `json:"Status"`
	Description string `json:"Description"`
}

func (t *ReplayBadRequest) Encode(w http.ResponseWriter) error {
	err := json.NewEncoder(w).Encode(&t)
	if !lib.LogOnError(err, "error: can't encode replyMq ReplayFetch") {
		return err
	}
	return nil
}

// SetBadRequest describe often used status
func (t *ReplayBadRequest) SetBadRequest(w http.ResponseWriter) {
	t.Status = http.StatusBadRequest
	t.Description = http.StatusText(http.StatusBadRequest)
	w.WriteHeader(http.StatusBadRequest)
}

func (t *ReplayBadRequest) SetNotFound(w http.ResponseWriter, desc string) {
	t.Status = http.StatusNotFound
	t.Description = desc
	w.WriteHeader(http.StatusNotFound)
}
