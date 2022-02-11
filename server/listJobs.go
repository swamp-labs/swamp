package main

import (
	"net/http"
)

type listJobsResponse struct {
	Function string `json:"function"`
}

// listJobs is the handler for GET - /api/v0/jobs
func listJobs(_ *http.Request) (interface{}, error) {
	return &listJobsResponse{"listJobsHandler"}, nil
}
