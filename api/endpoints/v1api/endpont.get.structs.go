package v1api

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/FetchTaskServer/lib"
	"net/http"
)

// ReqFetch get url
type ReqFetch struct {
	Method string `json:"Method"` // Method post or get
	Url    string `json:"Url"`    // URL
}

// Decode request
func (t *ReqFetch) Decode(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&t)
	if !lib.LogOnError(err, "error: can't decode request ReqFetch") {
		return err
	}
	return nil
}

func (t *ReqFetch) Validate() error {
	if t.Url == "" {
		return fmt.Errorf("url missed")
	}

	if t.Method == "" {
		return fmt.Errorf("method missed")
	}

	return nil
}

type ReplayBodyFetch struct {
	ID         int    `json:"ID"`
	StatusHttp int    `json:"StatusHttp"`
	Headers    string `json:"Headers"`
	Length     int    `json:"Length"` // answer Length
}

type ReplayFetch struct {
	Ok   int             `json:"Ok"`
	Body ReplayBodyFetch `json:"Body"`
}

func (t *ReplayFetch) Encode(w http.ResponseWriter) error {
	err := json.NewEncoder(w).Encode(&t)
	if !lib.LogOnError(err, "error: can't encode reply ReplayFetch") {
		return err
	}
	return nil
}
