// app.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// App export
type App struct {
	Router *mux.Router
}

func createJobHandler(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	json.Unmarshal([]byte(`{ "function": "createJobHandler" }`), &response)
	respondWithJSON(w, http.StatusOK, response)
}

func getJobHandler(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	json.Unmarshal([]byte(`{ "function": "getJobHandler" }`), &response)
	respondWithJSON(w, http.StatusOK, response)
}

func deleteJobHandler(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	json.Unmarshal([]byte(`{ "function": "deleteJobHandler" }`), &response)
	respondWithJSON(w, http.StatusOK, response)
}

func listJobsHandler(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	json.Unmarshal([]byte(`{ "function": "listJobsHandler" }`), &response)
	respondWithJSON(w, http.StatusOK, response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	json.Unmarshal([]byte(`{ "function": "healthHandler" }`), &response)
	respondWithJSON(w, http.StatusOK, response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (app *App) initialiseRoutes() {
	app.Router = mux.NewRouter()
	s := app.Router.PathPrefix("/api").Subrouter()
	s.HandleFunc("/v0/job", createJobHandler).Methods("POST")
	s.HandleFunc("/v0/job/{id:[0-9]+}", getJobHandler).Methods("GET")
	s.HandleFunc("/v0/job/{id:[0-9]+}", deleteJobHandler).Methods("DELETE")
	s.HandleFunc("/v0/jobs", listJobsHandler).Methods("GET")
	s.HandleFunc("/v0/health", healthHandler).Methods("GET")
}

func (app *App) run() {
	log.Println("Starting web server...")
	log.Fatal(http.ListenAndServe(":8080", app.Router))
}
