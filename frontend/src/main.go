package main

import (
	"cars/frontend/src/functions"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Serve static files for images
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Handle routes with corresponding handlers
	http.HandleFunc("/", functions.HomeHandler)
	http.HandleFunc("/cars", functions.CarHandler)
	http.HandleFunc("/contact", functions.ContactHandler)
	http.HandleFunc("/cars/", functions.CarDetailHandler)
	http.HandleFunc("/submit-contact", functions.SubmitContactHandler)
	http.HandleFunc("/filters", functions.AdvancedFiltersHandler)
	http.HandleFunc("/compare", functions.ComparisonsHandler) // Added compare handler

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
