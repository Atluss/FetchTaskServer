package v1api

import (
	"net/http"
)

const (
	V1ApiQueue = "v1"

	Ok                  = 200
	SyntaxError         = 400
	InternalServerError = 500
)

type HeadRequest interface {
	Request() // execute request
}

type ApiRun interface {
	Execute()  // запуск исполняющей функции в запросе
	Validate() // валидация данных
}

type ApiRequest struct {
	HeadRequest
	w *http.ResponseWriter
	r *http.Request
}
