package urlmap

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// URLHashMap is the main hashmap
type URLHashMap struct {
	rwmutex sync.RWMutex
	URLMap  map[string]string
}

type result struct {
	ShortURL string `bson:"shortened_url_code"`
	LongURL  string `bson:"original_url"`
}

// Map - The global map
var Map URLHashMap

// Init - Initialise the map
func Init() {
	Map.URLMap = make(map[string]string)
}

// Get - will fetch the short url from the hashmap
func (m *URLHashMap) Get(shortURL string) (string, bool) {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()

	val, exists := m.URLMap[shortURL]
	return val, exists
}

// Update - fetch from MongoDB and update the hashmap
func (m *URLHashMap) Update() {
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()

	connectionString := os.Getenv("MONGODB_URL")
	dbName := os.Getenv("DBNAME")
	collectionName := os.Getenv("COLLECTION_NAME")

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

	var collection *mongo.Collection

	fmt.Println("URLHashMap Update: Connected to MongoDB!")
	collection = client.Database(dbName).Collection(collectionName)
	
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	
	var results []result
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, r := range results {
		m.URLMap[r.ShortURL] = r.LongURL
	}
	
	fmt.Println(m.URLMap)
	fmt.Println("HashMap Update Complete")
}