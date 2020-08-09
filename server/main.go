package main

import (
	"encoding/json"
	"fmt"
	"github.com/couchbase/gocb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"strings"
)

type Movie struct {
	ID      string      `json:"id,omitempty"`
	Name    string      `json:"name,omitempty"`
	Genre   string      `json:"genre,omitempty"`
	Formats MovieFormat `json:"formats,omitempty"`
}

type MovieFormat struct {
	Digital bool `json:"digital,omitempty"`
	Bluray  bool `json:"bluray,omitempty"`
	Dvd     bool `json:"dvd,omitempty"`
}

var bucket *gocb.Bucket
var bucketName string

func ListEndpoint(w http.ResponseWriter, req *http.Request) {
	var movies []Movie
	query := gocb.NewN1qlQuery("SELECT `" + bucketName + "`.* FROM `" + bucketName + "`")
	query.Consistency(gocb.RequestPlus)
	rows, _ := bucket.ExecuteN1qlQuery(query, nil)
	var row Movie
	for rows.Next(&row) {
		movies = append(movies, row)
		row = Movie{}
	}
	if movies == nil {
		movies = make([]Movie, 0)
	}
	json.NewEncoder(w).Encode(movies)
}

func SearchEndpoint(w http.ResponseWriter, req *http.Request) {
	var movies []Movie
	params := mux.Vars(req)
	var n1qlParams []interface{}
	n1qlParams = append(n1qlParams, strings.ToLower(params["title"]))
	query := gocb.NewN1qlQuery("SELECT `" + bucketName + "`.* FROM `" + bucketName + "` WHERE LOWER(name) LIKE '%' || $1 || '%'")
	query.Consistency(gocb.RequestPlus)
	rows, _ := bucket.ExecuteN1qlQuery(query, n1qlParams)
	var row Movie
	for rows.Next(&row) {
		movies = append(movies, row)
		row = Movie{}
	}

	if movies == nil {
		movies = make([]Movie, 0)
	}
	json.NewEncoder(w).Encode(movies)
}

func CreateEndpoint(w http.ResponseWriter, req *http.Request) {
	var movie Movie
	_ = json.NewDecoder(req.Body).Decode(&movie)
	bucket.Insert(uuid.NewV4().String(), movie, 0)
	json.NewEncoder(w).Encode(movie)
}

func main() {
	fmt.Println("Starting server at http://localhost:12345...")
	cluster, _ := gocb.Connect("couchbase://172.17.0.2")
	err := cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "admin",
		Password: "admin12345",
	})
	if err != nil {
		panic(err)
	}
	bucketName = "go-restfull-sample"
	bucket, _ = cluster.OpenBucket(bucketName, "")
	router := mux.NewRouter()
	router.HandleFunc("/movies", ListEndpoint).Methods("GET")
	router.HandleFunc("/movies", CreateEndpoint).Methods("POST")
	router.HandleFunc("/search/{title}", SearchEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
