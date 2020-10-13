package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DSC-KIIT/divert/models"
	"github.com/lithammer/shortuuid"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	connectionString := os.Getenv("MONGODB_URL")
	dbName := os.Getenv("DBNAME")
	collName := os.Getenv("COLLNAME")

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	collection = client.Database(dbName).Collection(collName)
	fmt.Println("Collection instance created!")
}

// CreateShortenedURL for creating new shortened url
func CreateShortenedURL(w http.ResponseWriter, r *http.Request) {
	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// generate new unique id for url
	var reqURL models.URLShorten
	_ = json.NewDecoder(r.Body).Decode(&reqURL)
	id := shortuuid.New()
	reqURL.ShortenedURLCode = id
	reqURL.ClickCount = 0

	// insert in mongo
	result, err := collection.InsertOne(context.TODO(), reqURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	json.NewEncoder(w).Encode(result.InsertedID)
}
