package main

import (
	"net/http"
)

type getJobResponse struct {
	Function string `json:"function"`
}

// getJob is the handler for /api/v0/jobs/{id:[0-9]+}
func getJob(_ *http.Request) (interface{}, error) {
	return &getJobResponse{"getJobHandler"}, nil
}
