package urlmap

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/DSC-KIIT/divert/logger"
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
var lg logger.Logger

// Init - Initialise the map
func Init() {
	Map.URLMap = make(map[string]string)
	lg.Init()
}

// Get - will fetch the short url from the hashmap
func (m *URLHashMap) Get(shortURL string) (string, bool) {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()

	val, exists := m.URLMap[shortURL]
	return val, exists
}

// Fixes
func updateData(c chan<- int, data []result, m *URLHashMap) {
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()

	// Update the local data with new data in here
	fmt.Println("Updating data")
	for _, r := range data {
		m.URLMap[r.ShortURL] = r.LongURL
	}

	fmt.Println("Reads unlocked")
	c <- 10
	fmt.Println("Data updated")
}

func completed(c <-chan int, s chan<- int) {
	stat := <-c
	s <- stat
}

// Update - fetch from MongoDB and update the hashmap
func (m *URLHashMap) Update() {
	updateChannel := make(chan int, 1)
	defer close(updateChannel)
	status := make(chan int, 1)
	defer close(status)

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

	var collection *mongo.Collection

	lg.WriteInfo("URLHashMap Update: Connected to MongoDB!")
	collection = client.Database(dbName).Collection(collectionName)

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		lg.WriteError(err.Error())
	}

	var results []result
	if err = cursor.All(context.TODO(), &results); err != nil {
		lg.WriteError(err.Error())
	}

	updateData(updateChannel, results, m)
	completed(updateChannel, status)
	fmt.Println("Random update completion response : ", <-status)

	lg.WriteInfo("HashMap Update Complete")
}
