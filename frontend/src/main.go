package main

import (
	"cars/frontend/src/functions"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Define the routes and corresponding handlers
	http.HandleFunc("/", functions.HomeHandler)
	http.HandleFunc("/cars", functions.CarHandler)
	http.HandleFunc("/cars/", functions.CarDetailHandler)
	http.HandleFunc("/filters", functions.AdvancedFiltersHandler)
	http.HandleFunc("/compare", functions.ComparisonsHandler)

	// Serve static files like CSS, JS, images from the "static" directory
	staticDir := http.Dir(filepath.Join("frontend", "static"))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Start the HTTP server on port 8080
	log.Println("Starting server at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
