package app

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// CustomResponder - обертка для подключения к api handlers обработчиков типа http.HandlerFunc
type CustomResponder func(http.ResponseWriter, runtime.Producer)

func (c CustomResponder) WriteResponse(w http.ResponseWriter, p runtime.Producer) {
	c(w, p)
}
