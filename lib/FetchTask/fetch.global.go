package FetchTask

import (
	"fmt"
	"github.com/Atluss/FetchTaskServer/lib"
	"sync"
)

// global map of requests
var FetchElements map[string]FetchElement

// IsInElements check element in array
func IsInElements(token string) bool {

	if FetchElements == nil {
		FetchElements = map[string]FetchElement{}
	}

	if _, ok := FetchElements[token]; !ok {
		return false
	}

	return true
}

// AddToElements on global elements array
func AddToElements(obj *FetchElement) string {

	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	var id string
	for {
		id = lib.RandStringRunes(8)
		if !IsInElements(id) {
			obj.ID = id
			FetchElements[id] = *obj
			return id
		}
	}
}

// GetElementById get element by id
func GetElementById(id string) (*FetchElement, error) {

	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	if FetchElements == nil {
		return nil, fmt.Errorf("no element id: %s", id)
	}

	if el, ok := FetchElements[id]; !ok {
		return nil, fmt.Errorf("no element id: %s", id)
	} else {
		return &el, nil
	}

}
