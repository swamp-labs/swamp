package api

import (
	"encoding/json"
	"net/http"
)

// Handler functions can be used with the JsonHandler function to form a request handler
type Handler func(r *http.Request) (interface{}, error)

func handleError(w http.ResponseWriter, err error) {
	httpError, isHttpError := err.(*HTTPError)

	// Set error status code
	if isHttpError {
		w.WriteHeader(httpError.GetCode())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Return json body
	_ = json.NewEncoder(w).Encode(newErrorResponse(err))
}

// JsonHandler transform a function into a handler compatible with mux.
// It adds the error handler and the content type application/json header.
func JsonHandler(handler Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := handler(r)

		w.Header().Set("Content-Type", "application/json")

		// Handle errors
		if err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}
