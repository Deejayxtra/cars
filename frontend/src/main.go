package main

import (
	"cars/frontend/src/functions"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", functions.HomeHandler)
	http.HandleFunc("/cars", functions.CarHandler)
	http.HandleFunc("/cars/", functions.CarDetailHandler)
	http.HandleFunc("/filters", functions.AdvancedFiltersHandler)
	http.HandleFunc("/compare", functions.ComparisonsHandler)
	http.Handle("/src/static/", http.StripPrefix("/src/static/", http.FileServer(http.Dir("frontend/src/static"))))
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
