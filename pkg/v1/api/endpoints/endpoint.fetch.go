package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/FetchTaskServer/pkg/v1"
	"github.com/Atluss/FetchTaskServer/pkg/v1/FetchTask"
	"github.com/Atluss/FetchTaskServer/pkg/v1/api"
	"github.com/Atluss/FetchTaskServer/pkg/v1/config"
	"github.com/nats-io/go-nats"
	"log"
	"net/http"
	"sync"
	"time"
)

// NewEndPointV1Fetch constructor for /v1/test/{id}
func NewEndPointV1Fetch(set *config.Setup) (*v1fetch, error) {

	url := fmt.Sprintf("/%s/fetch", api.V1ApiQueue)
	if err := api.CheckEndPoint(api.V1ApiQueue, url); err != nil {
		return nil, err
	}

	endP := v1fetch{
		Setup: set,
		Url:   url,
	}
	return &endP, nil
}

type v1fetch struct {
	api.ApiRequest
	Setup *config.Setup
	Url   string
}

type v1fetchAnswer struct {
	request     *ReqFetch
	replyMq     *FetchTask.FetchElement
	replyClient *FetchTask.PublicElement
	badRequest  *api.ReplayBadRequest
}

// Request setup mux answer
func (obj *v1fetch) Request(w http.ResponseWriter, r *http.Request) {

	fErr := false
	req := &v1fetchAnswer{
		request:     &ReqFetch{},
		replyMq:     &FetchTask.FetchElement{},
		replyClient: &FetchTask.PublicElement{},
		badRequest:  &api.ReplayBadRequest{},
	}

	w.Header().Set("Content-Type", "application/json")

	// decode FetchTask body
	if err := req.request.Decode(r); err != nil {
		req.badRequest.SetBadRequest(w)
		v1.LogOnError(req.badRequest.Encode(w), "warning")
		return
	}

	log.Printf("Request: %s, params: %+v", obj.Url, req.request)

	// validate if decode is ok
	if err := req.request.Validate(); err != nil {
		fErr = true
		log.Printf("error: validate '%s' FetchTask %+v", err, req.request)
		req.badRequest.SetBadRequest(w)
		req.badRequest.Description = err.Error()
		v1.LogOnError(req.badRequest.Encode(w), "warning")
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {

		data, err := json.Marshal(&req.request)
		if err != nil {
			v1.LogOnError(err, fmt.Sprintf("warning: Problem with parsing FetchTask: %s", obj.Url))
			req.badRequest.SetBadRequest(w)
			fErr = true
		}

		if !fErr {
			msg, err := obj.Setup.Nats.Request(obj.Url, data, 31*time.Second)

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
			log.Printf("Request: %s done, id: %s", obj.Url, req.replyMq.ID)
			req.replyClient.SetFromElement(req.replyMq)
			w.WriteHeader(http.StatusOK)
			v1.LogOnError(req.replyClient.Encode(w), "warning")
		}

		wg.Done()
	}()
	wg.Wait()
}

// NatsQueue add new queue
func (obj *v1fetch) NatsQueue(m *nats.Msg) {

	answer := FetchTask.FetchElement{}

	params := ReqFetch{}
	err := json.Unmarshal(m.Data, &params)
	if !v1.LogOnError(err, "can't Unmarshal json") {
		return
	}

	if err := answer.GetOnline(params.Method, params.Url); err != nil {
		log.Println(err)
	}

	data, err := json.Marshal(&answer)
	if !v1.LogOnError(err, "can't Unmarshal json") {
		return
	}

	err = obj.Setup.Nats.Publish(m.Reply, data)
	v1.LogOnError(err, "warning")
}
