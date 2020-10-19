package router

import (
	"net/http"
	"time"

	"github.com/DSC-KIIT/divert/auth"
	"github.com/DSC-KIIT/divert/middleware"
	"github.com/DSC-KIIT/divert/urlmap"
	"github.com/gorilla/mux"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := mux.Vars(r)["shortURL"]
	longURL, exists := urlmap.Map.Get(shortURL)

	if exists {
		http.Redirect(w, r, longURL, 303)
		go middleware.IncrementClick(shortURL)
	} else {
		http.ServeFile(w, r, "public/404.html")
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

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/404.html")
}

// Router is exported and used in main.go
func Router() *mux.Router {
	middleware.Init()
	urlmap.Init()
	auth.Init()

	schedule(urlmap.Map.Update, 3*time.Minute)

	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET", "OPTIONS")
	router.HandleFunc("/{shortURL}", redirect).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/login", auth.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/createURL", middleware.CreateShortenedURL).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/getAllURL", middleware.GetAllURL).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/updateURL", middleware.UpdateShortURL).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deleteURL", middleware.DeleteURL).Methods("POST", "OPTIONS")

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	return router
}
