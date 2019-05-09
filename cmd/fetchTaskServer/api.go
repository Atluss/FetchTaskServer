package main

import (
	"fmt"
	"github.com/Atluss/FetchTaskServer/pkg/v1"
	"github.com/Atluss/FetchTaskServer/pkg/v1/api/controllers"
	"github.com/Atluss/FetchTaskServer/pkg/v1/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	settingPath := "settings.json"

	set := config.NewApiSetup(settingPath)

	log.Printf("Name: %s", set.Config.Name)
	log.Printf("Version: %s", set.Config.Version)
	log.Printf("Nats version: %s", set.Config.Nats.Version)
	log.Printf("Nats ReconnectedWait: %d", set.Config.Nats.ReconnectedWait)
	log.Printf("Nats host: %s", set.Config.Nats.Address[0].Host)
	log.Printf("Nats port: %s", set.Config.Nats.Address[0].Port)
	log.Printf("Nats address: %s", set.Config.Nats.Address[0].Address)
	log.Printf("Nats address(multi): %s", set.Config.GetNatsAddresses())

	// do something if user close program (close DB, or wait running query)
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Exit program...")
		os.Exit(1)
	}()

	// setup nats queue for test FetchTask
	err := controllers.NewV1Fetch(set)
	v1.LogOnError(err, "warning")

	err = controllers.NewV1Get(set)
	v1.LogOnError(err, "warning")

	err = controllers.NewV1List(set)
	v1.LogOnError(err, "warning")

	err = controllers.NewV1Delete(set)
	v1.LogOnError(err, "warning")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", set.Config.Port), set.Route))

}
