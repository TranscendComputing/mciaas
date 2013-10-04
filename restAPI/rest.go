package restAPI

import (
	"github.com/ant0ine/go-json-rest"
)

type IdResponse struct {
	Id string
}

func ProcessError(w *rest.ResponseWriter, err error, responseCode int) {
	if err != nil {
		rest.Error(w, err.Error(), responseCode)
	}
}
