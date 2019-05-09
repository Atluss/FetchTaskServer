package controllers

import (
	"github.com/Atluss/FetchTaskServer/pkg/v1"
	"github.com/Atluss/FetchTaskServer/pkg/v1/api"
	"github.com/Atluss/FetchTaskServer/pkg/v1/api/endpoints"
	"github.com/Atluss/FetchTaskServer/pkg/v1/config"
	"log"
)

// NewV1Fetch /v1/test/{id} register new Nats queue and frontend FetchTask
func NewV1Fetch(set *config.Setup) error {

	endP, err := endpoints.NewEndPointV1Fetch(set)
	if err != nil {
		return err
	}

	log.Printf("Setup endpoint: %s", endP.Url)

	// register queue for API and url
	_, err = set.Nats.QueueSubscribe(endP.Url, api.V1ApiQueue, endP.NatsQueue)
	if !v1.LogOnError(err, "Can't listen nats queue") {
		return err
	}

	// register frontend url
	set.Route.HandleFunc(endP.Url, endP.Request)

	api.AddEndPoint(api.V1ApiQueue, endP.Url)
	return nil
}
