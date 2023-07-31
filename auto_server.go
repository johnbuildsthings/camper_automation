package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"time"
)

// type Relay struct {
// 	ip          string
// 	description string
// 	enabled     bool
// }

const uri = "mongodb://localhost:27017"

// type BookRepository interface {
// 	FindById(ctx context.Context, id int) (*Book, error)
// }
// type bookRepository struct {
// 	client *mongo.Client
// }

// func (r *bookRepository) FindById(ctx context.Context, id int) (*Book, error) {
// 	var book Book
// 	err := r.client.DefaultDatabase().Collection("books").FindOne(ctx, bson.M{"_id": id}).Decode(&book)

// 	if err == mongo.ErrNoDocuments {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &book, nil
// }

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func dbSetup() *mongo.Client {
	// if client != nil {
	// 	return client
	// }
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	return client
}

func main() {
	// client := dbSetup()
	http.HandleFunc("/1/on", hello)
	http.HandleFunc("/1/off", headers)

	yfile, err := ioutil.ReadFile("items.yaml")

	if err != nil {
		panic(err)
	}

	data := make(map[interface{}]interface{})

	err2 := yaml.Unmarshal(yfile, &data)

	if err2 != nil {
		panic(err2)
	}

	for k, v := range data {

		fmt.Printf("%s -> %d\n", k, v)
	}
	// collection := client.Database("testing").Collection("numbers")
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// defer func() {
	// 	if err := client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// res, _ := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
	// id := res.InsertedID
	// fmt.Printf("id: %d", id)
	http.ListenAndServe(":8090", nil)
}
