package v1api

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/FetchTaskServer/lib"
	"net/http"
	"net/url"
)

// ReqFetch get url
type ReqFetch struct {
	Method string `json:"Method"` // Method post or get
	Url    string `json:"Url"`    // URL
}

// Decode FetchTask
func (t *ReqFetch) Decode(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&t)
	if !lib.LogOnError(err, "error: can't decode FetchTask ReqFetch") {
		return err
	}
	return nil
}

func (t *ReqFetch) Validate() error {
	if t.Url == "" {
		return fmt.Errorf("url missed")
	}

	_, err := url.ParseRequestURI(t.Url)
	if err != nil {
		return fmt.Errorf("url not correct")
	}

	if t.Method != http.MethodGet && t.Method != http.MethodPost {
		return fmt.Errorf("method missed")
	}

	return nil
}
