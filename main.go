package main

import (
	"context"
	"encoding/json"
	"log"
	"mongo-api/helper"
	"mongo-api/models"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Connection mongoDB with helper class
var (
	collection = helper.ConnectDB()
)

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Book array
	var books []models.Book

	// bson.M{}, we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		helper.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var book models.Book
		// character returns the memory address of the following variable.
		err := cur.Decode(&book) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}
		// add item to our array
		books = append(books, book)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	// encode similar to serialize process.
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	// set header
	w.Header().Set("Content-Type", "application/json")

	var book models.Book
	// we get params with mux
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&book)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&book)

	// insert out book model.
	result, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitive.objectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// prepare filter
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(deleteResult)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	// Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var book models.Book
	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&book)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"isbn", book.Isbn},
			{"title", book.Title},
			{"author", bson.D{
				{"firstname", book.Author.FirstName},
				{"lastname", book.Author.LastName},
			}},
		}},
	}

	if err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&book); err != nil {
		helper.GetError(err, w)
		return
	}
	book.ID = id
	json.NewEncoder(w).Encode(book)

}

func SetupRouters() (r *mux.Router) {

	r = mux.NewRouter()

	// arrange our route
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	return
}

func main() {
	r := SetupRouters()
	log.Println("Development Server Started On Port Number 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Println("error to start server", err.Error())
		return
	}
}
