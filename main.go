// package main

// import (
// 	"crypto/md5"
// 	"encoding/hex"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// type URL struct {
// 	ID           string    `json:"id"`
// 	OriginalURL  string    `json:"original_url"`
// 	ShortURL     string    `json:"short_url"`
// 	CreationDate time.Time `json:"creation_date"`
// }

// // create a databse
// var urlDB = make(map[string]URL)

// func generateUrlShortner(OriginalURL string) string {
// 	hasher := md5.New()
// 	hasher.Write([]byte(OriginalURL)) // it connverts the original url string into slice byte
// 	fmt.Println("hasher", hasher)
// 	data := hasher.Sum(nil)
// 	fmt.Println("Hasher data", data)
// 	hash := hex.EncodeToString(data)
// 	fmt.Println("Encode to string", hash)
// 	fmt.Println("Final string", hash[:8])
// 	return hash[:8]

// }
// func createURL(originalUrl string) string {
// 	shortUrl := generateUrlShortner(originalUrl)
// 	id := shortUrl //Use the short url as the ID for simplicity
// 	urlDB[id] = URL{
// 		ID:           id,
// 		OriginalURL:  originalUrl,
// 		ShortURL:     shortUrl,
// 		CreationDate: time.Now(),
// 	}
// 	return shortUrl
// }
// func getURL(id string) (URL, error) {
// 	url, ok := urlDB[id]
// 	if !ok {
// 		return URL{}, errors.New("URL not found")
// 	}
// 	return url, nil

// }
// func RootPageURL(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "This is the root page😁")
// }
// func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
// 	var data struct {
// 		URL string `json:"url"`
// 	}
// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}
// 	shortURL := createURL(data.URL)
// 	// fmt.Fprintf(w, shortURL)
// 	response := struct {
// 		ShortURL string `json:"short_url"`
// 	}{ShortURL: shortURL}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
// func RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Path[len("/redirect/"):]
// 	url, err := getURL(id)
// 	if err != nil {
// 		http.Error(w, "Invalid request", http.StatusNotFound)
// 	}
// 	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
// }
// func main() {
// 	// fmt.Println("URL Shortner started ......")
// 	// OriginalURL := "https://github.com/pradumkumar6"
// 	// generateUrlShortner(OriginalURL)
// 	// Register the handler function to handle all the requests to the root url("/")
// 	http.HandleFunc("/", RootPageURL)
// 	http.HandleFunc("/shorten", ShortURLHandler)
// 	http.HandleFunc("/redirect/", RedirectURLHandler)

// 	// setup the server
// 	fmt.Println("Server is running on port:8080")
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		fmt.Println("Error in server running", err)
// 	}

// }

package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

// create a database
var urlDB = make(map[string]URL)

func generateUrlShortner(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL)) // it converts the original url string into slice byte
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8]
}

func createURL(originalUrl string) string {
	shortUrl := generateUrlShortner(originalUrl)
	id := shortUrl // Use the short url as the ID for simplicity
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalUrl,
		ShortURL:     shortUrl,
		CreationDate: time.Now(),
	}
	return shortUrl
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func RootPageURL(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("static", "index.html"))
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	shortURL := createURL(data.URL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", RedirectURLHandler)

	fmt.Println("Server is running on port:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error in server running", err)
	}
}
