package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/FetchTaskServer/pkg/v1"
	"github.com/Atluss/FetchTaskServer/pkg/v1/FetchTask"
	"github.com/Atluss/FetchTaskServer/pkg/v1/api"
	"github.com/Atluss/FetchTaskServer/pkg/v1/config"
	"github.com/gorilla/mux"
	"github.com/nats-io/go-nats"
	"log"
	"net/http"
	"sync"
	"time"
)

// NewEndPointV1Delete constructor for /v1/delete/{id}
func NewEndPointV1Delete(set *config.Setup) (*v1delete, error) {

	url := fmt.Sprintf("/%s/delete/{id}", api.V1ApiQueue)

	if err := api.CheckEndPoint(api.V1ApiQueue, url); err != nil {
		return nil, err
	}

	endP := v1delete{
		Setup: set,
		Url:   url,
	}

	return &endP, nil

}

type v1delete struct {
	api.ApiRequest
	Setup *config.Setup
	Url   string
}

type v1deleteAnswer struct {
	replyMq     *FetchTask.FetchElement
	replyClient *FetchTask.PublicElement
	badRequest  *api.ReplayBadRequest
}

// Request setup mux answer
func (obj *v1delete) Request(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	fErr := false

	req := &v1deleteAnswer{
		replyMq:     &FetchTask.FetchElement{ID: vars["id"]},
		replyClient: &FetchTask.PublicElement{},
		badRequest:  &api.ReplayBadRequest{},
	}

	log.Printf("Request to delete element ID: %s", req.replyMq.ID)

	w.Header().Set("Content-Type", "application/json")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {

		data, err := json.Marshal(&req.replyMq)
		if err != nil || len(req.replyMq.ID) == 0 {
			v1.LogOnError(err, fmt.Sprintf("warning: Problem with parsing FetchTask: %s", obj.Url))
			req.badRequest.SetBadRequest(w)
			fErr = true
		}

		if !fErr {
			msg, err := obj.Setup.Nats.Request(obj.Url, data, 30*time.Second)

			if err == nil && msg != nil {
				err := json.Unmarshal(msg.Data, req.replyMq)
				if !v1.LogOnError(err, fmt.Sprintf("error: can't parse answer FetchTask %s", obj.Url)) {
					req.badRequest.SetBadRequest(w)
					fErr = true
				} else {
					if req.replyMq.Error != "" {
						req.badRequest.SetNotFound(w, req.replyMq.Error)
						fErr = true
					}
				}
			}
		}

		if fErr {
			v1.LogOnError(req.badRequest.Encode(w), "warning")
		} else {
			log.Printf("Request to detelte ID: %s done", req.replyMq.ID)
			req.badRequest.Status = http.StatusOK
			req.badRequest.Description = fmt.Sprintf("element id: %s deleted", req.replyMq.ID)

			w.WriteHeader(http.StatusOK)
			v1.LogOnError(json.NewEncoder(w).Encode(req.badRequest), "error: can't decode answer for list")
		}

		wg.Done()
	}()

	wg.Wait()
}

// NatsQueue add new queue
func (obj *v1delete) NatsQueue(m *nats.Msg) {

	answer := FetchTask.FetchElement{}

	err := json.Unmarshal(m.Data, &answer)
	if !v1.LogOnError(err, "can't Unmarshal params json") {
		return
	}

	err = FetchTask.DeleteFromList(answer.ID)
	if !v1.LogOnError(err, "warning") {
		answer.Error = err.Error()
	}

	data, err := json.Marshal(&answer)
	if !v1.LogOnError(err, "can't Unmarshal json") {
		return
	}

	err = obj.Setup.Nats.Publish(m.Reply, data)
	v1.LogOnError(err, "warning")
}
