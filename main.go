package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie // create a slice of movies of type Movie

/* Functions */

// get all movies in the struct
func getMovies(w http.ResponseWriter, r *http.Request)  {
	// setting the content type as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies) // encode the movies slice we have created above
}

// delete a movie (You will need to pas the id of the movie)
func deleteMovie(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // the id of the movie will be passed as a request param
	for index, item := range movies {
		if item.ID == params["id"] {
			// replace the movie at index with movie at (index +1) and others following
			// this effectively deletes the movie at index matching id passed in params
			movies = append(movies[:index], movies[index+1:]...) 
			break
		}
	}
	// return the remaining movies 
	json.NewEncoder(w).Encode(movies)
}

// get a movie by id
func getMovie(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // collect the id of the movie as a param
	for _, item := range movies {
		if item.ID == params["id"] {
			// return the movie with the id
			json.NewEncoder(w).Encode(item)
			break
		}
	}
}

// create A movie
func createMovie(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	// define a variable called movie of type movie
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000)) // generate a random number for the id
	movies = append(movies, movie) // append the movie created to the movies slice
	json.NewEncoder(w).Encode(movie) // returns the movie that has been created
}

// update a movie
func updateMovie(w http.ResponseWriter, r *http.Request)  {
	// set json content type
	w.Header().Set("Content-Type", "application/json")

	// params
	params:= mux.Vars(r)

	// range over the movies in the movies slice
	// delete the movie with id
	// add the movie we have created
	for index, item := range movies {
		if item.ID == params["id"] {
			// deleting the movie of provided ID
			movies = append(movies[:index], movies[index+1:]...)
			// creating  a new movie with random ID
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"] // using the same ID for the updated movie instance
			movies = append(movies, movie) // append the movie to the movies slice
			json.NewEncoder(w).Encode(movie) // return the newly created movie
			return
		}
	}
}


/******* Main Function Starts Here ********/
func main() {
	// assign newRouter from the mux gorrilla library to r
	r := mux.NewRouter()

	// fetch a default movie instance when movies API is successfully called
	movies = append(movies, Movie{ID: "1", Isbn: "48227", Title: "Rhapsody of Realities Movie", Director: &Director{FirstName: "Clinton", LastName: "Felix"}})
	movies = append(movies, Movie{ID: "2", Isbn: "45455", Title: "Healing School Movie", Director: &Director{FirstName: "Pst Deola", LastName: "Philips"}})

	// set Handler functions
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}