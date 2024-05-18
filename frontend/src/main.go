package main

import (
	"cars/frontend/src/functions"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Serve static files for frontend
	http.Handle("/src/static/", http.StripPrefix("/src/static/", http.FileServer(http.Dir("frontend/src/static"))))

	// Serve static files for images
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../cars/api/src/img"))))

	// Handle other routes
	http.HandleFunc("/", functions.HomeHandler)
	http.HandleFunc("/cars", functions.CarHandler)
	http.HandleFunc("/cars/", functions.CarDetailHandler) // Use the new CarDetailHandler
	http.HandleFunc("/filters", functions.AdvancedFiltersHandler)
	http.HandleFunc("/compare", functions.ComparisonsHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
