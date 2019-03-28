package FetchTask

import (
	"encoding/json"
	"github.com/Atluss/FetchTaskServer/lib"
	"net/http"
	"time"
)

type FetchElement struct {
	ID      string      `json:"ID"`
	Status  int         `json:"StatusHttp"`
	Headers http.Header `json:"Headers"`
	Length  int64       `json:"Length"` // answer Length
}

// Encode encode answer
func (t *FetchElement) Encode(w http.ResponseWriter) error {
	err := json.NewEncoder(w).Encode(&t)
	if !lib.LogOnError(err, "error: can't encode reply ReplayFetch") {
		return err
	}
	return nil
}

func (t *FetchElement) Get(method, url string) error {

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	t.Status = resp.StatusCode
	t.Length = resp.ContentLength
	t.Headers = resp.Header

	AddToElements(t)

	return nil
}
