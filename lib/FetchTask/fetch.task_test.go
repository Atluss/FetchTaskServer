package FetchTask

import (
	"log"
	"testing"
)

func TestFetchElement_Get(t *testing.T) {

	ft := FetchElement{}

	if err := ft.Get("GET", "https://ya.ru"); err != nil {
		log.Println(err)
	} else {
		log.Printf("%+v", ft)
	}

}
