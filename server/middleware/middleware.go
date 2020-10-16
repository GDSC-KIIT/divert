package middleware

import (
	"context"
	"encoding/json"
	"fmt"
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
var lg logger.Logger

type response struct {
	Status  string `json:"status"`
	Message interface{} `json:"message"`
}

// Init - Connect to the DB
func Init() {
	lg.Init()

	connectionString := os.Getenv("MONGODB_URL")
	dbName := os.Getenv("DBNAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		lg.WriteError(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		lg.WriteError(err.Error())
	}
	defer cancel()

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		lg.WriteError(err.Error())
	}

	lg.WriteInfo("Middleware Init: Connected to MongoDB!")
	collection = client.Database(dbName).Collection(collectionName)
	lg.WriteInfo("Middleware Init: Collection instance created!")
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

	lg.WriteInfo(fmt.Sprintf("Creating Short URL: %v to Long URL: %v", reqURL.ShortenedURLCode, reqURL.OriginalURL))

	// check if the same short url already exists
	var searchResult bson.M
	err := collection.FindOne(context.TODO(), bson.D{{"shortened_url_code", reqURL.ShortenedURLCode}}).Decode(&searchResult)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// short url is unique -> insert in db
			result, err := collection.InsertOne(context.TODO(), reqURL)
			if err != nil {
				lg.WriteError(err.Error())
			}

			json.NewEncoder(w).Encode(response{"okay", result.InsertedID})
		}
	} else {
		json.NewEncoder(w).Encode(response{"error", "Same Short URL already exists"})
	}
}

// GetAllURL - Get all the URLs in the DB
func GetAllURL(w http.ResponseWriter, r *http.Request) {
	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		lg.WriteError(err.Error())
	}
	defer cursor.Close(context.TODO())

	var urls []bson.M

	for cursor.Next(context.TODO()) {
		var url bson.M
		if err = cursor.Decode(&url); err != nil {
			lg.WriteError(err.Error())
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

	lg.WriteInfo(fmt.Sprintf("Update URL: %v with new Long URL %v", reqURL.ShortenedURLCode, reqURL.OriginalURL))

	filter := bson.D{{"_id", reqURL.ID}}
	replacement := bson.D{{"$set", bson.D{{"original_url", reqURL.OriginalURL}}}}

	var replaced bson.M
	err := collection.FindOneAndUpdate(context.TODO(), filter, replacement).Decode(&replaced)
	if err != nil {
		lg.WriteError(err.Error())
	}

	json.NewEncoder(w).Encode(response{"okay", "Updated successfully"})
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
		lg.WriteError(err.Error())
	}

	json.NewEncoder(w).Encode(response{"okay", fmt.Sprintf("deleted %v documents", result.DeletedCount)})
}

// IncrementClick - Increments the click counter on the db
func IncrementClick(shortURL string) {
	filter := bson.D{{"shortened_url_code", shortURL}}
	update := bson.D{{"$inc", bson.D{{"click_count", 1}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		lg.WriteError(err.Error())
	}
}