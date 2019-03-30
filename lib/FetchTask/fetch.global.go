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
func GetElementById(id string) (FetchElement, error) {

	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	el := FetchElement{}
	var ok bool

	if FetchElements == nil {
		return el, fmt.Errorf("no element id: %s", id)
	}

	if el, ok = FetchElements[id]; !ok {
		return el, fmt.Errorf("no element id: %s", id)
	} else {
		return el, nil
	}

}

func GetListElement() []FetchElement {

	ret := []FetchElement{}

	for _, v := range FetchElements {
		ret = append(ret, v)
	}

	return ret
}

// pages start at 1, can't be 0 or less.
//func GetDataPage(page, perPage int) map[string]interface{} {
//	start := (page - 1) * perPage
//	stop := start + perPage
//
//	if start > len(FetchElements) {
//		return nil
//	}
//
//	if stop > len(FetchElements) {
//		stop = len(FetchElements)
//	}
//
//	return FetchElements[start:stop]
//}
