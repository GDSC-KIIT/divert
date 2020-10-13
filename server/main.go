package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DSC-KIIT/divert/router"
)

func main() {
	r := router.Router()
	port := os.Getenv("PORT")
	fmt.Println("Starting server on the port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
