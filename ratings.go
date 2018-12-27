package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Movie contains movie information
type Movie struct {
	ID     int     `json:"id"`
	Rating float32 `json:"rating"`
}

// NewRating creates a movie with new rating
func (movie *Movie) NewRating() (newMovie *Movie) {
	rating := rand.Float32() * 10

	return &Movie{
		ID:     movie.ID,
		Rating: rating,
	}
}

// ToJSON converts movie object to json
func (movie *Movie) ToJSON() string {
	data, err := json.Marshal(movie)
	if err != nil {
		// Should not happen in production, so panic to fail fast and warn developer
		panic(err)
	}

	return string(data)
}

// ReadMovies reads a movies file and creates a map of movies
func ReadMovies(filename string) []Movie {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		// Should not happen in production, so panic to fail fast and warn developer
		panic(err)
	}

	var movieList []Movie
	err = json.Unmarshal(data, &movieList)
	if err != nil {
		// Should not happen in production, so panic to fail fast and warn developer
		panic(err)
	}
	return movieList
}

// SSE writes Server-Sent Events to an HTTP client.
type SSE struct{}

var messages = make(chan string)

func (s *SSE) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	f, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Error: Cannot stream", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	cn, ok := rw.(http.CloseNotifier)
	if !ok {
		http.Error(rw, "Error: Cannot stream", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-cn.CloseNotify():
			log.Println("Done: Closed connection")
			return
		case msg := <-messages:
			fmt.Fprintf(rw, "{\"data\": %s}\r\n", msg)
			f.Flush()
		}
	}
}

func main() {
	movies := ReadMovies("./data/movie_ratings.json")
	moviesTotal := len(movies)

	http.Handle("/ratings", &SSE{})
	log.Println("Started: Serving stream...")
	go func() {
		for {
			movieNo := rand.Intn(moviesTotal)
			messages <- movies[movieNo].NewRating().ToJSON()
			delay := time.Duration(rand.Intn(500)) * time.Millisecond
			time.Sleep(delay)
		}
	}()

	log.Fatal(http.ListenAndServe(":3000", nil))
}
