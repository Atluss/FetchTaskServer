package FetchTask

import (
	"log"
	"testing"
)

func TestAddToElements(t *testing.T) {

	ft := FetchElement{}

	if err := ft.Get("GET", "https://ya.ru"); err != nil {
		log.Println(err)
	} else {
		//log.Println(AddToElements(&ft))

		log.Printf("%+v", FetchElements)
	}

}

func TestIsInElements(t *testing.T) {

	id := "sew"

	log.Printf("%t", IsInElements(id))

}
