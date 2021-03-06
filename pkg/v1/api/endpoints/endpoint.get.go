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

// NewEndPointV1Get constructor for /v1/test/{id}
func NewEndPointV1Get(set *config.Setup) (*v1get, error) {

	url := fmt.Sprintf("/%s/get/{id}", api.V1ApiQueue)

	if err := api.CheckEndPoint(api.V1ApiQueue, url); err != nil {
		return nil, err
	}

	endP := v1get{
		Setup: set,
		Url:   url,
	}

	return &endP, nil

}

type v1get struct {
	api.ApiRequest
	Setup *config.Setup
	Url   string
}

type v1getAnswer struct {
	replyMq     *FetchTask.FetchElement
	replyClient *FetchTask.PublicElement
	badRequest  *api.ReplayBadRequest
}

// Request setup mux answer
func (obj *v1get) Request(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	fErr := false

	req := &v1getAnswer{
		replyMq:     &FetchTask.FetchElement{ID: vars["id"]},
		replyClient: &FetchTask.PublicElement{},
		badRequest:  &api.ReplayBadRequest{},
	}

	log.Printf("Request: %s, id: %s", obj.Url, req.replyMq.ID)

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
			log.Printf("Request: %s done for ID: %s", obj.Url, req.replyMq.ID)
			req.replyClient.SetFromElement(req.replyMq)
			w.WriteHeader(http.StatusOK)
			v1.LogOnError(req.replyClient.Encode(w), "warning")
		}

		wg.Done()
	}()

	wg.Wait()
}

// NatsQueue add new queue
func (obj *v1get) NatsQueue(m *nats.Msg) {

	answer := FetchTask.FetchElement{}

	err := json.Unmarshal(m.Data, &answer)
	if !v1.LogOnError(err, "can't Unmarshal params json") {
		return
	}

	answer, err = FetchTask.GetElementById(answer.ID)
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
