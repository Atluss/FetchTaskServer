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
func NewEndPointV1Get(set *config.Setup) (*v1get, error) {

	url := fmt.Sprintf("/%s/fetch", V1ApiQueue)

	if err := api.CheckEndPoint(V1ApiQueue, url); err != nil {
		return nil, err
	}

	endP := v1get{
		Setup: set,
		Url:   url,
	}

	return &endP, nil

}

type v1get struct {
	ApiRequest
	Setup *config.Setup
	Url   string
}

type v1getAnswer struct {
	request    *ReqFetch
	reply      *FetchTask.FetchElement
	badRequest *ReplayBadRequest
}

// Request setup mux answer
func (obj *v1get) Request(w http.ResponseWriter, r *http.Request) {

	fErr := false
	req := &v1getAnswer{
		request:    &ReqFetch{},
		reply:      &FetchTask.FetchElement{},
		badRequest: &ReplayBadRequest{},
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
				err := json.Unmarshal(msg.Data, req.reply)

				log.Printf("%+v", req.reply)

				if !lib.LogOnError(err, fmt.Sprintf("error: can't parse answer FetchTask %s", obj.Url)) {
					req.badRequest.SetBadRequest(w)
					fErr = true
				}

			}
		}

		if fErr {
			lib.LogOnError(req.badRequest.Encode(w), "warning")
		} else {
			log.Printf("Answer: %+v: for FetchTask: %+v", req.reply, req.request)
			lib.LogOnError(req.reply.Encode(w), "warning")
		}

		wg.Done()

	}()
	wg.Wait()

}

// NatsQueue add new queue
func (obj *v1get) NatsQueue(m *nats.Msg) {

	answer := FetchTask.FetchElement{}

	params := ReqFetch{}
	err := json.Unmarshal(m.Data, &params)

	log.Printf("%+v", params)
	if !lib.LogOnError(err, "can't unpasrse json") {
		return
	}

	if err := answer.Get(params.Method, params.Url); err != nil {
		log.Println(err)
	}

	data, err := json.Marshal(&answer)
	if !lib.LogOnError(err, "can't unpasrse json") {
		return
	}

	log.Println("Replying to ", m.Reply)

	err = obj.Setup.Nats.Publish(m.Reply, data)
	lib.LogOnError(err, "warning")
}
