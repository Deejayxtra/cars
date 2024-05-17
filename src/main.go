package main

import (
    "fmt"
	"cars/functions"
	"net/http"
)


func main() {
    http.HandleFunc("/", functions.HomeHandler)
    http.HandleFunc("/cars", functions.CarHandler)
    http.HandleFunc("/filters", functions.AdvancedFiltersHandler)
    http.HandleFunc("/compare", functions.ComparisonsHandler)

    // Serve static files
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))

    // Start the HTTP server
    
    port := 8080
    fmt.Printf("Server is running on port %d\n", port)
    fmt.Println(http.ListenAndServe(":8080", nil))
}


