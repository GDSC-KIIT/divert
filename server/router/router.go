package router

import (
	"github.com/DSC-KIIT/divert/middleware"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/createUrl", middleware.CreateShortenedURL).Methods("POST", "OPTIONS")

	return router
}
