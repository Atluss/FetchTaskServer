package controllers

import (
	"github.com/Atluss/FetchTaskServer/pkg/v1"
	"github.com/Atluss/FetchTaskServer/pkg/v1/api"
	"github.com/Atluss/FetchTaskServer/pkg/v1/api/endpoints"
	"github.com/Atluss/FetchTaskServer/pkg/v1/config"
	"log"
)

// NewV1Get return endpoint new fetch element
func NewV1Get(set *config.Setup) error {

	endP, err := endpoints.NewEndPointV1Get(set)
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
