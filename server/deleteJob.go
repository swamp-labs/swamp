package main

import (
	"net/http"
)

type deleteJobResponse struct {
	Function string `json:"function"`
}

// deleteJob is the handler for DELETE - /api/v0/jobs/{id:[0-9]+}
func deleteJob(_ *http.Request) (interface{}, error) {
	return &deleteJobResponse{"getJobHandler"}, nil
}
