package router

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DSC-KIIT/divert/auth"
	"github.com/DSC-KIIT/divert/middleware"
	"github.com/DSC-KIIT/divert/urlmap"
	"github.com/gorilla/mux"
)

func secure(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self';base-uri 'self';block-all-mixed-content;font-src 'self' https: data:;frame-ancestors 'self';img-src 'self' data:;object-src 'none';script-src 'self';script-src-attr 'none';style-src 'self' https: 'unsafe-inline';upgrade-insecure-requests")
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		w.Header().Set("Expect-CT", "max-age=0")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Strict-Transport-Security", "max-age=15552000; includeSubDomains")
		w.Header().Set("X-Download-Options", "noopen")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		h.ServeHTTP(w, r)
	})
}

func easteregg(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.Proto, r.URL)

		if r.URL.Path == "/omgdsc" {
			fmt.Fprint(w, "Hello from DSC-KIIT, you found an easter egg 🐰🥚")
			return
		}

		h.ServeHTTP(w, r)
	})
}

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
	router.Use(secure)
	router.Use(easteregg)
	fs := http.FileServer(http.Dir("public/"))

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
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
