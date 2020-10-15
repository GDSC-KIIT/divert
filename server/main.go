package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DSC-KIIT/divert/router"
	// "github.com/DSC-KIIT/divert/logger"
)

func main() {
	os.Setenv("PORT", "3000")
	//os.Setenv("MONGODB_URL", "mongodb+srv://admin:adminpassword@cluster0.q3enw.mongodb.net/divert?retryWrites=true&w=majority")
	os.Setenv("MONGODB_URL", "mongodb://localhost:27017")
	os.Setenv("DBNAME", "divert")
	os.Setenv("COLLECTION_NAME", "urls")


	r := router.Router()
	port := os.Getenv("PORT")
	fmt.Println("Starting server on the port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, r))

	// 	var l logger.Logger
	// 	l.Init();

	// 	l.WriteInfo("I am info");
	// 	l.WriteWarning("I am warning");
	// 	l.WriteError("I am error");
}
