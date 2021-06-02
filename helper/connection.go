package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB: This is helper function to connect mongoDB
// If you want to export your function. You must to start case function name.
// Otherwise you won't see your function when you import that on other class.

func ConnectDB() *mongo.Collection {
	// Set client Options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error to make connection with mongoDB.")
		return nil
	}
	fmt.Println("Connected to MongoDB!")
	collection := client.Database("go_rest_api").Collection("books")
	return collection
}

// ErrorResponse: This is error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError: This is helper function to prepare error model
func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}
	message, _ := json.Marshal(response)
	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
