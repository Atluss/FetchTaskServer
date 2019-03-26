package api

import (
	"log"
	"testing"
)

func TestAddEndPoint(t *testing.T) {

	queue := "v1"
	url := "/test"

	AddEndPoint(queue, url)

	for i, k := range EndPoints {
		log.Println("=============")
		log.Printf("Queue: %s", i)

		for _, nc := range k {
			log.Printf("%s", nc.url)
		}
		log.Println("=============")
	}
}
