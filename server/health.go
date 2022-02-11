package main

import (
	"net/http"
)

type healthResponse struct {
	Function string `json:"function"`
}

// health is the handler for /api/v0/health
func health(_ *http.Request) (interface{}, error) {
	return &healthResponse{"healthHandler"}, nil
}
