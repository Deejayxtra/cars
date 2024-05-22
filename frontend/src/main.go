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

	// Handle other routes
	http.HandleFunc("/", functions.HomeHandler)
	http.HandleFunc("/cars", functions.CarHandler)
	http.HandleFunc("/cars/", functions.CarDetailHandler) // Use the new CarDetailHandler
	http.HandleFunc("/filters", functions.AdvancedFiltersHandler)
	http.HandleFunc("/compare", functions.ComparisonsHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
