package main

import (
	"net/http"
)

type createJobResponse struct {
	Function string `json:"function"`
}

// createJob is the handler for POST - /api/v0/jobs
func createJob(_ *http.Request) (interface{}, error) {
	return &createJobResponse{"createJobHandler"}, nil
}
