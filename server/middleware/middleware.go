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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/DSC-KIIT/divert/logger"
)

var collection *mongo.Collection

// Init - Connect to the DB
func Init() {
	var l logger.Logger
	l.Init()

	connectionString := os.Getenv("MONGODB_URL")
	dbName := os.Getenv("DBNAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		l.WriteError(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		l.WriteError(err.Error())
	}
	defer cancel()

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		l.WriteError(err.Error())
	}

	fmt.Println("Connected to MongoDB!")
	collection = client.Database(dbName).Collection(collectionName)
	fmt.Println("Collection instance created!")
}

// Index - Get Request of the Index page
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "DSCKIIT - divert, Backend API")
}

// CreateShortenedURL for creating new shortened url
func CreateShortenedURL(w http.ResponseWriter, r *http.Request) {
	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// set short and long url from r.Body to the reqURL obj
	var reqURL models.URLShorten
	_ = json.NewDecoder(r.Body).Decode(&reqURL)
	reqURL.ClickCount = 0

	fmt.Printf("Create: %v", reqURL)

	// insert in mongo
	result, err := collection.InsertOne(context.TODO(), reqURL)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(result.InsertedID)
}

// GetAllURL - Get all the URLs in the DB
func GetAllURL(w http.ResponseWriter, r *http.Request) {
	// TODO: make sure you cannot insert duplicates
	var l logger.Logger
	l.Init()

	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		l.WriteError(err.Error())
	}
	defer cursor.Close(context.TODO())

	var urls []bson.M

	for cursor.Next(context.TODO()) {
		var url bson.M
		if err = cursor.Decode(&url); err != nil {
			l.WriteError(err.Error())
		}
		urls = append(urls, url)
	}

	json.NewEncoder(w).Encode(urls)
}

// UpdateShortURL - Update a particular short URL
func UpdateShortURL(w http.ResponseWriter, r *http.Request) {
	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var reqURL models.URLShorten
	_ = json.NewDecoder(r.Body).Decode(&reqURL)

	fmt.Printf("Update URL Body: %v", reqURL)

	filter := bson.D{{"_id", reqURL.ID}}
	replacement := bson.D{{"$set", bson.D{{"original_url", reqURL.OriginalURL}}}}

	var replaced bson.M
	err := collection.FindOneAndUpdate(context.TODO(), filter, replacement).Decode(&replaced)
	if err != nil {
		log.Fatal(err)
	}

	type response struct {
		Status string `json:"status"`
	}

	resp := response{"okay"}
	json.NewEncoder(w).Encode(resp)
}

// DeleteURL - Delete a particular URL
func DeleteURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var reqURL models.URLShorten
	_ = json.NewDecoder(r.Body).Decode(&reqURL)

	filter := bson.D{{"_id", reqURL.ID}}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	type response struct {
		Status string `json:"status"`
	}

	resp := response{fmt.Sprintf("deleted %v documents", result.DeletedCount)}

	json.NewEncoder(w).Encode(resp)
}
