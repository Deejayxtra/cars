package main

import (
	"cars/frontend/src/functions"
	"log"
	"fmt"
	"net/http"
	"path/filepath"
)

func main() {
	// Serve static files for images
 	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

 	// Handle other routes
 	// Define the routes and corresponding handlers
 	http.HandleFunc("/", functions.HomeHandler)
	http.HandleFunc("/cars", functions.CarHandler)
	http.HandleFunc("/contact", functions.ContactHandler)
	http.HandleFunc("/cars/", functions.CarDetailHandler) // Use the new CarDetailHandler
	http.HandleFunc("/submit-contact", functions.SubmitContactHandler)
	http.HandleFunc("/filters", functions.AdvancedFiltersHandler)
	http.HandleFunc("/compare", functions.ComparisonsHandler)

 	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	// Serve static files like CSS, JS, images from the "static" directory
 	staticDir := http.Dir(filepath.Join("frontend", "static"))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

 	// Start the HTTP server on port 8080
 	log.Println("Starting server at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
