package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DSC-KIIT/divert/middleware"
	"github.com/DSC-KIIT/divert/urlmap"
	"github.com/gorilla/mux"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := mux.Vars(r)["shortURL"]
	longURL, exists := urlmap.Map.Get(shortURL)
	
	if exists {
		http.Redirect(w, r, longURL, 303)
	} else {
		fmt.Fprintf(w, "DSCKIIT Divert - 404 Page Not found")
	}
}

func schedule(f func(), delay time.Duration) {
	go func() {
		for {
			f()
			time.Sleep(delay)
		}
	}()
}

// Router is exported and used in main.go
func Router() *mux.Router {
	middleware.Init()
	urlmap.Init()

	schedule(urlmap.Map.Update, 3*time.Minute)

	router := mux.NewRouter()

	router.HandleFunc("/{shortURL}", redirect).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/createURL", middleware.CreateShortenedURL).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/getAllURL", middleware.GetAllURL).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/updateURL", middleware.UpdateShortURL).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deleteURL", middleware.DeleteURL).Methods("POST", "OPTIONS")
	router.HandleFunc("/api", middleware.Index).Methods("GET", "OPTIONS")

	return router
}
