package web

import (
	"aurum/config"
	"aurum/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// Starts the REST web API
func StartServer(config *config.Config, db db.Connection) {

	endpoints := Endpoints{
		conn:   db,
		config: config,
	}

	// Router
	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/signup", endpoints.signup).Methods("POST")
	api.HandleFunc("/login", endpoints.login).Methods("POST")

	/// Create the server
	srv := &http.Server{
		Handler: router,
		Addr:    config.WebURL,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting up web server ...")
	log.Fatal(srv.ListenAndServe())
}
