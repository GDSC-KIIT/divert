package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DSC-KIIT/divert/logger"
	"github.com/DSC-KIIT/divert/models"
	jwt "github.com/dgrijalva/jwt-go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var authCollection *mongo.Collection
var lg logger.Logger

type response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

// Init - Connect to DB and initialise the auth collection
func Init() {
	lg.Init()

	connectionString := os.Getenv("MONGODB_URL")
	dbName := os.Getenv("DBNAME")
	collectionName := os.Getenv("AUTH_COLLECTION_NAME")

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

	authCollection = client.Database(dbName).Collection(collectionName)
}

func generateJWT(username string) string {
	signingKey := []byte(os.Getenv("JWT_SIGNING_KEY"))

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		lg.WriteError(err.Error())
	}

	return tokenString
}

// IsValidToken - checks if the passed auth token is valid or expired
func IsValidToken(authToken string) bool {
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid Token")
		}

		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})

	if err != nil {
		lg.WriteWarning(err.Error())
		return false
	}

	return token.Valid
}

func isValidUser(username string, password string) bool {
	var searchResult models.AuthModel
	err := authCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&searchResult)

	if err != nil {
		lg.WriteWarning(err.Error())
		return false
	}

	return password == searchResult.Password
}

// Login - route for login
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req models.AuthModel
	_ = json.NewDecoder(r.Body).Decode(&req)

	token := r.Header.Get("x-auth-token")

	if token == "" {
		if isValidUser(req.Username, req.Password) {
			newToken := generateJWT(req.Username)
			resp := response{Status: "ok", Message: "New Token Issued", Token: newToken}
			json.NewEncoder(w).Encode(resp)
		} else {
			resp := response{Status: "error", Message: "Incorrect Username Password"}
			json.NewEncoder(w).Encode(resp)
		}
	} else {
		if IsValidToken(token) {
			resp := response{Status: "ok", Message: "Token is Valid"}
			json.NewEncoder(w).Encode(resp)
		} else {
			if isValidUser(req.Username, req.Password) {
				newToken := generateJWT(req.Username)
				resp := response{Status: "ok", Message: "New Token Issued", Token: newToken}
				json.NewEncoder(w).Encode(resp)
			} else {
				resp := response{Status: "error", Message: "Incorrect Username Password"}
				json.NewEncoder(w).Encode(resp)
			}
		}
	}
}
