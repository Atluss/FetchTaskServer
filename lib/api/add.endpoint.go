package api

import (
	"fmt"
)

type endPoint struct {
	url string
}

var EndPoints map[string][]endPoint

func AddEndPoint(queue, url string) {

	if EndPoints == nil {
		EndPoints = map[string][]endPoint{}
	}

	_, ok := EndPoints[queue]
	if !ok {
		EndPoints[queue] = []endPoint{}
	}

	EndPoints[queue] = append(EndPoints[queue], endPoint{url: url})
}

func CheckEndPoint(queue, url string) error {

	if EndPoints == nil {
		return nil
	}

	_, ok := EndPoints[queue]
	if !ok {
		return nil
	}

	for _, nc := range EndPoints[queue] {
		if nc.url == url {
			return fmt.Errorf("endpoint: %s already set", url)
		}
	}

	return nil

}
