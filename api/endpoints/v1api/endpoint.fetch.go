package v1api

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/FetchTaskServer/lib"
	"github.com/Atluss/FetchTaskServer/lib/FetchTask"
	"github.com/Atluss/FetchTaskServer/lib/api"
	"github.com/Atluss/FetchTaskServer/lib/config"
	"github.com/nats-io/go-nats"
	"log"
	"net/http"
	"sync"
	"time"
)

// NewEndPointV1Test constructor for /v1/test/{id}
func NewEndPointV1Fetch(set *config.Setup) (*v1fetch, error) {

	url := fmt.Sprintf("/%s/fetch", V1ApiQueue)

	if err := api.CheckEndPoint(V1ApiQueue, url); err != nil {
		return nil, err
	}

	endP := v1fetch{
		Setup: set,
		Url:   url,
	}

	return &endP, nil

}

type v1fetch struct {
	ApiRequest
	Setup *config.Setup
	Url   string
}

type v1fetchAnswer struct {
	request     *ReqFetch
	replyMq     *FetchTask.FetchElement
	replyClient *FetchTask.PublicElement
	badRequest  *ReplayBadRequest
}

// Request setup mux answer
func (obj *v1fetch) Request(w http.ResponseWriter, r *http.Request) {

	fErr := false
	req := &v1fetchAnswer{
		request:     &ReqFetch{},
		replyMq:     &FetchTask.FetchElement{},
		replyClient: &FetchTask.PublicElement{},
		badRequest:  &ReplayBadRequest{},
	}

	w.Header().Set("Content-Type", "application/json")

	// decode FetchTask body
	if err := req.request.Decode(r); err != nil {
		fErr = true
	}

	// validate if decode is ok
	if !fErr {
		if err := req.request.Validate(); err != nil {
			fErr = true
			log.Printf("error: validate '%s' FetchTask %+v", err, req.request)
		}
	}

	if fErr {

		req.badRequest.SetBadRequest(w)
		lib.LogOnError(req.badRequest.Encode(w), "warning")

		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {

		data, err := json.Marshal(&req.request)
		if err != nil {
			lib.LogOnError(err, fmt.Sprintf("warning: Problem with parsing FetchTask: %s", obj.Url))
			req.badRequest.SetBadRequest(w)
			fErr = true
		}

		if !fErr {
			msg, err := obj.Setup.Nats.Request(obj.Url, data, 31*time.Second)

			if err == nil && msg != nil {
				err := json.Unmarshal(msg.Data, req.replyMq)
				log.Printf("%+v", req.replyMq)

				if !lib.LogOnError(err, fmt.Sprintf("error: can't parse answer FetchTask %s", obj.Url)) {
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
			lib.LogOnError(req.badRequest.Encode(w), "warning")
		} else {
			log.Printf("Answer: %+v: for FetchTask: %+v", req.replyMq, req.request)
			req.replyClient.SetFromElement(req.replyMq)
			lib.LogOnError(req.replyClient.Encode(w), "warning")
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

	log.Printf("%+v", params)
	if !lib.LogOnError(err, "can't Unmarshal json") {
		return
	}

	if err := answer.GetOnline(params.Method, params.Url); err != nil {
		log.Println(err)
	}

	data, err := json.Marshal(&answer)
	if !lib.LogOnError(err, "can't Unmarshal json") {
		return
	}

	log.Println("Replying to ", m.Reply)

	err = obj.Setup.Nats.Publish(m.Reply, data)
	lib.LogOnError(err, "warning")
}
