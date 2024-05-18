package functions

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

// Define structs to represent the data retrieved from the API
type Specifications struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain   string `json:"drivetrain"`
}

type Car struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	ManufacturerID int            `json:"manufacturerId"`
	CategoryID     int            `json:"categoryId"`
	Year           int            `json:"year"`
	Specifications Specifications `json:"specifications"`
	Image          string         `json:"image"`
}

type CarDetail struct {
	Car          Car
	Manufacturer Manufacturer
}

type Manufacturer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type CarWithManufacturer struct {
	Car              Car    // Car details
	ManufacturerName string // Manufacturer name
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
	// Make a request to the Cars API to fetch all cars
	resp, err := http.Get("http://localhost:3000/api/models")
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

	// Decode JSON response to get all cars
	var cars []Car
	if err := json.NewDecoder(resp.Body).Decode(&cars); err != nil {
		fmt.Printf("Failed to decode car data: %v\n", err)
		http.Error(w, "Failed to decode car data. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Fetch manufacturers for filtering dropdown
	manufacturers, err := fetchManufacturers()
	if err != nil {
		fmt.Printf("Failed to fetch manufacturers: %v\n", err)
		http.Error(w, "Failed to fetch manufacturers. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Render HTML template with fetched data (all cars and manufacturers)
	tmpl, err := template.ParseFiles("templates/cars.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template with all cars and manufacturers data
	err = tmpl.Execute(w, struct {
		Cars          []Car
		Manufacturers []Manufacturer
	}{
		Cars:          cars,
		Manufacturers: manufacturers,
	})
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// fetchManufacturers fetches manufacturers data from the API
func fetchManufacturers() ([]Manufacturer, error) {
	// Make a request to the Manufacturers API
	resp, err := http.Get("http://localhost:3000/api/manufacturers")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the API responded with an error status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned an error: %s", resp.Status)
	}

	// Decode JSON response
	var manufacturers []Manufacturer
	if err := json.NewDecoder(resp.Body).Decode(&manufacturers); err != nil {
		return nil, err
	}

	return manufacturers, nil
}

func CarDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the car ID from the URL path
	carID := r.URL.Path[len("/cars/"):]

	log.Printf("Fetching details for car ID: %s", carID)

	// Fetch car details
	carResp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/models/%s", carID))
	if err != nil {
		log.Printf("Error fetching car details: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer carResp.Body.Close()

	if carResp.StatusCode != http.StatusOK {
		log.Printf("Car API returned non-OK status: %s", carResp.Status)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	carData, err := ioutil.ReadAll(carResp.Body)
	if err != nil {
		log.Printf("Error reading car details: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var car Car
	if err := json.Unmarshal(carData, &car); err != nil {
		log.Printf("Error unmarshaling car details: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch manufacturer details
	manufResp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/manufacturers/%d", car.ManufacturerID))
	if err != nil {
		log.Printf("Error fetching manufacturer details: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer manufResp.Body.Close()

	if manufResp.StatusCode != http.StatusOK {
		log.Printf("Manufacturer API returned non-OK status: %s", manufResp.Status)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	manufData, err := ioutil.ReadAll(manufResp.Body)
	if err != nil {
		log.Printf("Error reading manufacturer details: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var manufacturer Manufacturer
	if err := json.Unmarshal(manufData, &manufacturer); err != nil {
		log.Printf("Error unmarshaling manufacturer details: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Fetched car details: %+v", car)
	log.Printf("Fetched manufacturer details: %+v", manufacturer)

	// Combine car and manufacturer details
	carDetail := CarDetail{
		Car:          car,
		Manufacturer: manufacturer,
	}

	// Parse the template
	tmplPath := filepath.Join("templates", "carDetail.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing template %s: %v", tmplPath, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template
	if err := tmpl.Execute(w, carDetail); err != nil {
		log.Printf("Error executing template %s: %v", tmplPath, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func AdvancedFiltersHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/filters.html"))
	tmpl.Execute(w, nil)
}

func ComparisonsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/compare.html"))
	tmpl.Execute(w, nil)
}

func GetManufacturerName(manufacturerID int) (string, error) {
	// Make a request to the API endpoint to fetch manufacturers
	resp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/manufacturers/%d", manufacturerID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned an error: %s", resp.Status)
	}

	// Decode JSON response
	var manufacturer struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&manufacturer); err != nil {
		return "", err
	}

	return manufacturer.Name, nil
}
