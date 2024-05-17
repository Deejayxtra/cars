package functions

import (
	"fmt"
	"encoding/json"
	"html/template"
	"net/http"
)

// Define structs to represent the data retrieved from the API
type Car struct {
	ID           int    `json:"id"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
	Year         int    `json:"year"`
	// Add other fields as needed
}

// CarsResponse represents the structure of the API response
type CarsResponse struct {
    Cars []Car `json:"cars"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/home.html")
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, nil)
    if err != nil {
        fmt.Printf("Error executing template: %v\n", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}

// CarHandler handles requests to render the cars template
func CarHandler(w http.ResponseWriter, r *http.Request) {
    // Make a request to the Cars API
    resp, err := http.Get("http://localhost:3000/api/")
    if err != nil {
        fmt.Printf("Failed to fetch car data: %v\n", err)
        http.Error(w, "Failed to fetch car data. Please try again later.", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

	// Check if the API responded with an error status code
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("API returned an error: %s\n", resp.Status)
        http.Error(w, "Failed to fetch car data. Please try again later.", http.StatusInternalServerError)
        return
    }

    // Decode JSON response
    var carsResponse CarsResponse
    err = json.NewDecoder(resp.Body).Decode(&carsResponse)
    if err != nil {
        fmt.Printf("Failed to decode car data: %v\n", err)
        http.Error(w, "Failed to decode car data. Please try again later.", http.StatusInternalServerError)
        return
    }

    // Render HTML template with fetched data
    tmpl, err := template.ParseFiles("templates/cars.html")
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, carsResponse.Cars)
    if err != nil {
        fmt.Printf("Error executing template: %v\n", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}


// Handler to render advanced filters template
func AdvancedFiltersHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/filters.html"))
    tmpl.Execute(w, nil)
}

// Handler to render comparisons template
func ComparisonsHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/compare.html"))
    tmpl.Execute(w, nil)
}