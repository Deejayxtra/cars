package functions

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
	// Define channels for data and error handling
	carChan := make(chan []Car)
	manufChan := make(chan []Manufacturer)
	errChan := make(chan error)

	// Get query parameters
	search := r.URL.Query().Get("search")
	manufacturerID := r.URL.Query().Get("manufacturer")

	// Fetch car data asynchronously
	go fetchCarsAsync(carChan, errChan)

	// Fetch manufacturers data asynchronously
	go fetchManufacturersAsync(manufChan, errChan)

	// Wait for the data or error from channels
	var cars []Car
	var manufacturers []Manufacturer

	for i := 0; i < 2; i++ {
		select {
		case fetchedCars := <-carChan:
			cars = fetchedCars
		case fetchedManufacturers := <-manufChan:
			manufacturers = fetchedManufacturers
		case err := <-errChan:
			log.Printf("Error occurred: %v", err)
			http.Error(w, "Sorry, something went wrong on our end. Please try again later.", http.StatusInternalServerError)
			return
		}
	}

	// Filter cars based on search query and manufacturer ID
	var filteredCars []Car
	for _, car := range cars {
		if (search == "" || strings.Contains(strings.ToLower(car.Name), strings.ToLower(search))) &&
			(manufacturerID == "" || strconv.Itoa(car.ManufacturerID) == manufacturerID) {
			filteredCars = append(filteredCars, car)
		}
	}

	// Render HTML template with fetched data (filtered cars and manufacturers)
	tmpl, err := template.ParseFiles("templates/cars.html")
	if err != nil {
		log.Printf("Error parsing template: %v\n", err)
		http.Error(w, "Sorry, something went wrong on our end. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Execute the template with filtered cars and manufacturers data
	err = tmpl.Execute(w, struct {
		Cars          []Car
		Manufacturers []Manufacturer
	}{
		Cars:          filteredCars,
		Manufacturers: manufacturers,
	})
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Sorry, something went wrong on our end. Please try again later.", http.StatusInternalServerError)
		return
	}
}

// fetchCarsAsync fetches car data asynchronously
func fetchCarsAsync(carChan chan []Car, errChan chan error) {
	resp, err := http.Get("http://localhost:3000/api/models")
	if err != nil {
		errChan <- fmt.Errorf("failed to fetch car data: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("API returned an error: %s", resp.Status)
		return
	}

	var cars []Car
	if err := json.NewDecoder(resp.Body).Decode(&cars); err != nil {
		errChan <- fmt.Errorf("failed to decode car data: %v", err)
		return
	}

	carChan <- cars
}

// fetchManufacturersAsync fetches manufacturer data asynchronously
func fetchManufacturersAsync(manufChan chan []Manufacturer, errChan chan error) {
	resp, err := http.Get("http://localhost:3000/api/manufacturers")
	if err != nil {
		errChan <- fmt.Errorf("failed to fetch manufacturers: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("Manufacturer API returned an error: %s", resp.Status)
		return
	}

	var manufacturers []Manufacturer
	if err := json.NewDecoder(resp.Body).Decode(&manufacturers); err != nil {
		errChan <- fmt.Errorf("failed to decode manufacturers data: %v", err)
		return
	}

	manufChan <- manufacturers
}

func CarDetailHandler(w http.ResponseWriter, r *http.Request) {
	carID := r.URL.Path[len("/cars/"):]
	carChan := make(chan Car)
	manufChan := make(chan Manufacturer)
	errChan := make(chan error)

	// Fetch car and manufacturer details asynchronously
	go fetchCarDetailAsync(carID, carChan, errChan)
	go fetchManufacturerDetailAsync(carID, manufChan, errChan)

	var car Car
	var manufacturer Manufacturer

	for i := 0; i < 2; i++ {
		select {
		case fetchedCar := <-carChan:
			car = fetchedCar
		case fetchedManufacturer := <-manufChan:
			manufacturer = fetchedManufacturer
		case err := <-errChan:
			log.Printf("Error occurred: %v", err)
			http.Error(w, "Sorry, something went wrong on our end. Please try again later.", http.StatusInternalServerError)
			return
		}
	}

	carDetail := CarDetail{
		Car:          car,
		Manufacturer: manufacturer,
	}

	tmplPath := filepath.Join("templates", "carDetail.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing template %s: %v", tmplPath, err)
		http.Error(w, "Sorry, something went wrong on our end. Please try again later.", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, carDetail); err != nil {
		log.Printf("Error executing template %s: %v", tmplPath, err)
		http.Error(w, "Sorry, something went wrong on our end. Please try again later.", http.StatusInternalServerError)
	}
}

// fetchCarDetailAsync fetches car detail asynchronously
func fetchCarDetailAsync(carID string, carChan chan Car, errChan chan error) {
	carAPIURL := fmt.Sprintf("http://localhost:3000/api/models/%s", carID)
	resp, err := http.Get(carAPIURL)
	if err != nil {
		errChan <- fmt.Errorf("failed to fetch car details: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("Car API returned an error: %s", resp.Status)
		return
	}

	var car Car
	if err := json.NewDecoder(resp.Body).Decode(&car); err != nil {
		errChan <- fmt.Errorf("failed to decode car details: %v", err)
		return
	}

	carChan <- car
}

// fetchManufacturerDetailAsync fetches manufacturer detail asynchronously
func fetchManufacturerDetailAsync(carID string, manufChan chan Manufacturer, errChan chan error) {
	carAPIURL := fmt.Sprintf("http://localhost:3000/api/models/%s", carID)
	carResp, err := http.Get(carAPIURL)
	if err != nil {
		errChan <- fmt.Errorf("failed to fetch car details: %v", err)
		return
	}
	defer carResp.Body.Close()

	if carResp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("Car API returned an error: %s", carResp.Status)
		return
	}

	var car Car
	if err := json.NewDecoder(carResp.Body).Decode(&car); err != nil {
		errChan <- fmt.Errorf("failed to decode car details: %v", err)
		return
	}

	manufAPIURL := fmt.Sprintf("http://localhost:3000/api/manufacturers/%d", car.ManufacturerID)
	resp, err := http.Get(manufAPIURL)
	if err != nil {
		errChan <- fmt.Errorf("failed to fetch manufacturer details: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("Manufacturer API returned an error: %s", resp.Status)
		return
	}

	var manufacturer Manufacturer
	if err := json.NewDecoder(resp.Body).Decode(&manufacturer); err != nil {
		errChan <- fmt.Errorf("failed to decode manufacturer details: %v", err)
		return
	}

	manufChan <- manufacturer
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
	resp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/manufacturers/%d", manufacturerID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned an error: %s", resp.Status)
	}

	var manufacturer struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&manufacturer); err != nil {
		return "", err
	}

	return manufacturer.Name, nil
}
