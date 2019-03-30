package apiV1

import (
	"github.com/Atluss/FetchTaskServer/api/endpoints/v1api"
	"github.com/Atluss/FetchTaskServer/lib"
	"github.com/Atluss/FetchTaskServer/lib/api"
	"github.com/Atluss/FetchTaskServer/lib/config"
	"log"
)

func NewV1Delete(set *config.Setup) error {

	endP, err := v1api.NewEndPointV1Delete(set)
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
