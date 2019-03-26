package v1api

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/FetchTaskServer/lib"
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

type User struct {
	Id   string
	Name string
}

type v1getAnswer struct {
	request *ReqFetch
	reply   *ReplayFetch
}

// Request setup mux answer
func (obj *v1get) Request(w http.ResponseWriter, r *http.Request) {

	fErr := false
	req := &v1getAnswer{
		request: &ReqFetch{},
		reply:   &ReplayFetch{},
	}

	w.Header().Set("Content-Type", "application/json")

	// decode request body
	if err := req.request.Decode(r); err != nil {
		fErr = true
	}

	// validate if decode is ok
	if !fErr {
		if err := req.request.Validate(); err != nil {
			fErr = true
			log.Printf("error: validate '%s' request %+v", err, req.request)
		}
	}

	if fErr {
		req.reply.Ok = SyntaxError
		w.WriteHeader(400)
		lib.LogOnError(req.reply.Encode(w), "warning")
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {

		data, err := json.Marshal(&req.request)
		if err != nil {
			req.reply.Ok = SyntaxError
			lib.LogOnError(err, fmt.Sprintf("warning: Problem with parsing request: %s", obj.Url))
			w.WriteHeader(400)
			fErr = true
		}

		if !fErr {
			msg, err := obj.Setup.Nats.Request(obj.Url, data, 100*time.Millisecond)

			if err == nil && msg != nil {
				err := json.Unmarshal(msg.Data, &req.reply)
				if !lib.LogOnError(err, fmt.Sprintf("error: can't parse answer request %s", obj.Url)) {
					req.reply.Ok = SyntaxError
				} else {
					req.reply.Ok = Ok
				}
			}
		}

		log.Printf("Answer: %+v: for request: %+v", req.reply, req.request)
		lib.LogOnError(req.reply.Encode(w), "warning")
		wg.Done()

	}()
	wg.Wait()

}

// NatsQueue add new queue
func (obj *v1get) NatsQueue(m *nats.Msg) {

	answer := ReplayFetch{
		Body: ReplayBodyFetch{
			ID:         1,
			StatusHttp: 200,
			Headers:    "===",
			Length:     123,
		},
	}

	params := ReqFetch{}
	err := json.Unmarshal(m.Data, &params)

	log.Printf("%+v", params)

	if !lib.LogOnError(err, "can't unpasrse json") {
		return
	}

	data, err := json.Marshal(&answer)
	if !lib.LogOnError(err, "can't unpasrse json") {
		return
	}

	log.Println("Replying to ", m.Reply)

	err = obj.Setup.Nats.Publish(m.Reply, data)
	lib.LogOnError(err, "warning")
}
