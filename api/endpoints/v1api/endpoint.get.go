package v1api

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/Go-Nats-Api-Example/lib"
	"github.com/Atluss/Go-Nats-Api-Example/lib/api"
	"github.com/Atluss/Go-Nats-Api-Example/lib/config"
	"github.com/gorilla/mux"
	"github.com/nats-io/go-nats"
	"log"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Id   string
	Name string
}

// NewEndPointV1Test constructor for /v1/test/{id}
func NewEndPointV1Get(set *config.Setup) (*v1get, error) {

	url := fmt.Sprintf("/%s/test/{id}", V1ApiQueue)

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
	Setup *config.Setup
	Url   string
}

// Request setup mux answer
func (obj *v1get) Request(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	myUser := User{Id: vars["id"]}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {

		data, err := json.Marshal(&myUser)
		if err != nil || len(myUser.Id) == 0 {
			lib.LogOnError(err, "warning: Problem with parsing the user Id: %s")
			w.WriteHeader(500)
			return
		}

		msg, err := obj.Setup.Nats.Request(obj.Url, data, 100*time.Millisecond)
		if err == nil && msg != nil {

			myUserWithName := User{}

			err := json.Unmarshal(msg.Data, &myUserWithName)
			log.Printf("==== %+v", myUserWithName)

			if err == nil {
				myUser = myUserWithName
			}

			w.Header().Set("Content-Type", "application/json")

			lib.LogOnError(json.NewEncoder(w).Encode(myUserWithName), "warning")
		}
		wg.Done()

	}()
	wg.Wait()

}

// NatsQueue add new queue
func (obj *v1get) NatsQueue(m *nats.Msg) {

	users := map[string]string{
		"1": "Ilya",
		"2": "Yana",
		"3": "Olga",
		"4": "Shokin",
	}

	myUser := User{}
	err := json.Unmarshal(m.Data, &myUser)

	if err != nil {
		log.Println(err)
		return
	}

	myUser.Name = users[myUser.Id]

	data, err := json.Marshal(&myUser)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Replying to ", m.Reply)

	err = obj.Setup.Nats.Publish(m.Reply, data)
	lib.LogOnError(err, "warning")
}
