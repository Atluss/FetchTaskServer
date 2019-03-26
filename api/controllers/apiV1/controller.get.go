package apiV1

import (
	"github.com/Atluss/Go-Nats-Api-Example/api/endpoints/v1api"
	"github.com/Atluss/Go-Nats-Api-Example/lib"
	"github.com/Atluss/Go-Nats-Api-Example/lib/api"
	"github.com/Atluss/Go-Nats-Api-Example/lib/config"
	"log"
)

// NewV1Test /v1/test/{id} register new Nats queue and frontend request
func NewV1Test(set *config.Setup) error {

	endP, err := v1api.NewEndPointV1Get(set)
	if err != nil {
		return err
	}

	log.Printf("Setup endpoint: %s", endP.Url)

	// register queue for API and url
	_, err = set.Nats.QueueSubscribe(endP.Url, v1api.V1ApiQueue, endP.NatsQueue)
	if !lib.LogOnError(err, "Can't listen nats queue") {
		return err
	}

	// register frontend url
	set.Route.HandleFunc(endP.Url, endP.Request)

	api.AddEndPoint(v1api.V1ApiQueue, endP.Url)
	return nil
}
