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

var templates = template.Must(template.ParseFiles(
	"templates/home.html",
	"templates/contact.html",
	"templates/cars.html",
	"templates/carDetail.html",
	"templates/compare.html", // Added compare.html
))

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

// ComparisonsHandler handles the comparison of cars
func ComparisonsHandler(w http.ResponseWriter, r *http.Request) {
	carIDs := r.URL.Query()["ids"] // Assume car IDs are passed as query parameters like /compare?ids=1,2,3
	if len(carIDs) == 0 {
		http.Error(w, "No car IDs provided for comparison", http.StatusBadRequest)
		return
	}

	carChan := make(chan []Car)
	errChan := make(chan error)

	go fetchCarsByIdsAsync(carIDs, carChan, errChan)

	var cars []Car
	select {
	case fetchedCars := <-carChan:
		cars = fetchedCars
	case err := <-errChan:
		log.Printf("Error occurred: %v", err)
		http.Error(w, "Sorry, something went wrong on our end. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Render compare.html template with fetched car details
	tmplPath := filepath.Join("templates", "compare.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing template %s: %v", tmplPath, err)
		http.Error(w, "Sorry, something went wrong on our end. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Create a struct to pass to the template
	data := struct {
		Cars []Car
	}{
		Cars: cars,
	}

	// Execute the template with the structured data
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template %s: %v", tmplPath, err)
		http.Error(w, "Sorry, something went wrong on our end. Please try again later.", http.StatusInternalServerError)
	}
}

// fetchCarsByIdsAsync fetches car data asynchronously by car IDs
func fetchCarsByIdsAsync(carIDs []string, carChan chan []Car, errChan chan error) {
	var cars []Car

	for _, id := range carIDs {
		carAPIURL := fmt.Sprintf("http://localhost:3000/api/models/%s", id)
		resp, err := http.Get(carAPIURL)
		if err != nil {
			errChan <- fmt.Errorf("failed to fetch car details for ID %s: %v", id, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errChan <- fmt.Errorf("Car API returned an error for ID %s: %s", id, resp.Status)
			return
		}

		var car Car
		if err := json.NewDecoder(resp.Body).Decode(&car); err != nil {
			errChan <- fmt.Errorf("failed to decode car details for ID %s: %v", id, err)
			return
		}

		cars = append(cars, car)
	}

	carChan <- cars
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

func AdvancedFiltersHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/filters.html"))
	tmpl.Execute(w, nil)
}

// contactHandler serves the contact form page
func ContactHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "contact.html", nil); err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}
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

	// Render HTML template with fetched car details
	tmpl, err := template.ParseFiles("templates/carDetail.html")
	if err != nil {
		log.Printf("Error parsing template: %v\n", err)
		http.Error(w, "Sorry, something went wrong on our end. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Create CarDetail struct to pass to the template
	carDetail := CarDetail{
		Car:          car,
		Manufacturer: manufacturer,
	}

	// Execute the template with the car details
	err = tmpl.Execute(w, carDetail)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Sorry, something went wrong on our end. Please try again later.", http.StatusInternalServerError)
		return
	}
}

// fetchCarDetailAsync fetches car detail data asynchronously
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

// fetchManufacturerDetailAsync fetches manufacturer detail data asynchronously
func fetchManufacturerDetailAsync(carID string, manufChan chan Manufacturer, errChan chan error) {
	// Fetch car details first to get manufacturer ID
	carAPIURL := fmt.Sprintf("http://localhost:3000/api/models/%s", carID)
	resp, err := http.Get(carAPIURL)
	if err != nil {
		errChan <- fmt.Errorf("failed to fetch car details for manufacturer: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("Car API returned an error for manufacturer: %s", resp.Status)
		return
	}

	var car Car
	if err := json.NewDecoder(resp.Body).Decode(&car); err != nil {
		errChan <- fmt.Errorf("failed to decode car details for manufacturer: %v", err)
		return
	}

	// Fetch manufacturer details
	manufAPIURL := fmt.Sprintf("http://localhost:3000/api/manufacturers/%d", car.ManufacturerID)
	resp, err = http.Get(manufAPIURL)
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

// submitContactHandler handles the form submission
func SubmitContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	message := r.FormValue("message")

	// Handle the form submission (e.g., send an email, save to a database, etc.)
	log.Printf("Received contact form submission: Name: %s, Email: %s, Message: %s", name, email, message)

	// Pass the confirmation message to the template
	data := struct {
		Message string
	}{
		Message: "Thank you for your message. We will get back to you shortly.",
	}

	if err := templates.ExecuteTemplate(w, "contact.html", data); err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}
}
