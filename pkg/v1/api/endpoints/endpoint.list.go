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

// NewEndPointV1List constructor for /v1/test/{id}
func NewEndPointV1List(set *config.Setup) (*v1list, error) {

	url := fmt.Sprintf("/%s/list", api.V1ApiQueue)

	if err := api.CheckEndPoint(api.V1ApiQueue, url); err != nil {
		return nil, err
	}

	endP := v1list{
		Setup: set,
		Url:   url,
	}

	return &endP, nil

}

type v1list struct {
	api.ApiRequest
	Setup *config.Setup
	Url   string
}

type v1listAnswer struct {
	replyMq     *[]FetchTask.FetchElement
	replyClient []FetchTask.PublicElement
	badRequest  *api.ReplayBadRequest
}

// Request setup mux answer
func (obj *v1list) Request(w http.ResponseWriter, r *http.Request) {

	fErr := false

	req := &v1listAnswer{
		replyMq:     &[]FetchTask.FetchElement{},
		replyClient: []FetchTask.PublicElement{},
		badRequest:  &api.ReplayBadRequest{},
	}

	log.Printf("Request: %s", obj.Url)

	w.Header().Set("Content-Type", "application/json")

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {

		var data []byte

		if !fErr {
			msg, err := obj.Setup.Nats.Request(obj.Url, data, 30*time.Second)

			if err == nil && msg != nil {
				err := json.Unmarshal(msg.Data, req.replyMq)
				if !v1.LogOnError(err, fmt.Sprintf("error: can't parse answer FetchTask %s", obj.Url)) {
					req.badRequest.SetBadRequest(w)
					fErr = true
				} else {
					if len(*req.replyMq) == 0 {
						req.badRequest.SetNotFound(w, "no elements")
						fErr = true
					}
				}
			}
		}

		if fErr {
			v1.LogOnError(req.badRequest.Encode(w), "warning")
		} else {
			for _, v := range *req.replyMq {
				req.replyClient = append(req.replyClient, FetchTask.PublicElement{
					ID:      v.ID,
					Status:  v.Status,
					Headers: v.Headers,
					Length:  v.Length})
			}

			log.Printf("Request list done elements: %d", len(req.replyClient))
			w.WriteHeader(http.StatusOK)
			v1.LogOnError(json.NewEncoder(w).Encode(req.replyClient), "error: can't decode answer for list")
		}

		wg.Done()
	}()
	wg.Wait()
}

// NatsQueue add new queue
func (obj *v1list) NatsQueue(m *nats.Msg) {

	answer := FetchTask.GetListElement()
	data, err := json.Marshal(&answer)
	if !v1.LogOnError(err, "can't Unmarshal json") {
		return
	}
	err = obj.Setup.Nats.Publish(m.Reply, data)
	v1.LogOnError(err, "warning")

}
