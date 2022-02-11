// app.go
package main

import (
	"github.com/swamp-labs/swamp/server/api"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// App export
type App struct {
	Router *mux.Router
}

func (app *App) initialiseRoutes() {
	app.Router = mux.NewRouter()
	s := app.Router.PathPrefix("/api").Subrouter()
	s.HandleFunc("/v0/jobs", api.JsonHandler(createJob)).Methods("POST")
	s.HandleFunc("/v0/jobs/{id:[0-9]+}", api.JsonHandler(getJob)).Methods("GET")
	s.HandleFunc("/v0/jobs/{id:[0-9]+}", api.JsonHandler(deleteJob)).Methods("DELETE")
	s.HandleFunc("/v0/jobs", api.JsonHandler(listJobs)).Methods("GET")
	s.HandleFunc("/v0/health", api.JsonHandler(health)).Methods("GET")
}

func (app *App) run() {
	log.Println("Starting web server...")
	log.Fatal(http.ListenAndServe(":8080", app.Router))
}
