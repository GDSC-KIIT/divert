package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DSC-KIIT/divert/router"
)

func main() {
	// For local development purposes
	// os.Setenv("PORT", "3000")
	// os.Setenv("MONGODB_URL", "mongodb://localhost:27017")
	// os.Setenv("DBNAME", "divert")
	// os.Setenv("COLLECTION_NAME", "urls")
	// os.Setenv("AUTH_COLLECTION_NAME", "authinfo")
	// os.Setenv("JWT_SIGNING_KEY", "junaidrahim")

	port := os.Getenv("PORT")

	r := router.Router()
	fmt.Println("DSCKIIT Divert Backend Service - Starting server on the port " + port)
	fmt.Println("Logs written on stderr and divert.log file")

	log.Fatal(http.ListenAndServe(":"+port, r))
}
