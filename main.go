package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	store = make(map[string]struct {
		Filename string
		Data     []byte
	}) // In-memory key-value store
	mu sync.Mutex // Mutex for concurrent access
)

func main() {
	http.HandleFunc("/health", healthHandler)

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download/", downloadHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the file temporarily
	var data []byte
	data, err = io.ReadAll(file) // Use io.ReadAll instead of ioutil.ReadAll
	if err != nil {
		http.Error(w, "Unable to read file data", http.StatusInternalServerError)
		return
	}

	// Generate a random code for the file
	code := generateRandomCode()

	// Store the file data temporarily
	mu.Lock()
	// store[code] = string(data)
	store[code] = struct {
		Filename string
		Data     []byte
	}{Filename: header.Filename, Data: data}
	mu.Unlock()

	w.Write([]byte("File uploaded successfully. Download code: " + code))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/download/"):]

	mu.Lock()
	entry, exists := store[code]
	mu.Unlock()

	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set Content-Disposition header to preserve original filename
	w.Header().Set("Content-Disposition", "attachment; filename="+entry.Filename)
	w.Write(entry.Data)
}

// Create a new random generator with a seed based on the current time
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomCode() string {
	return strconv.Itoa(rng.Intn(100000)) // Use the new random generator
}
