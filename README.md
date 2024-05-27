# Car Viewer Application

This is a Car Viewer web application built using Go (Golang). It allows users to view car details, filter cars by manufacturer, and search for cars by name. The application fetches data from an API and displays it in a user-friendly interface.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Setup and Installation](#setup-and-installation)
- [Running the Application](#running-the-application)
- [Endpoints](#endpoints)
- [License](#license)

## Features

- Home page with a welcome message and a link to view cars.
- Display a list of cars fetched from an API.
- Search for cars by name.
- Filter cars by manufacturer.
- View detailed information about a car, including specifications and manufacturer details.

## Project Structure
├── api
│   ├── Makefile
│   ├── README.md
│   ├── data.json
│   ├── img
│   │   └── carsImages.jpg..
│   ├── main.js
│   ├── node_modules

├── frontend
│ ├── src
│ │ ├── functions
│ │ │ └── carHandler.go
│ │ ├── main.go
│ │ └── templates
│ │ ├── carDetail.html
│ │ ├── cars.html
│ │ └── home.html
│ └── static
│ ├── homePageImg
│ │ └── car-reviews.jpg
│ └── styles.css
├── README.md



## Setup and Installation

1. Clone the repository:
    ```sh
    git clone https://gitea.koodsisu.fi/olufemiekundayo/cars
    cd cars/frontend
    ```

2. Ensure you have Go installed. If not, download and install it from [https://golang.org/dl/](https://golang.org/dl/).

3. Install dependencies (if any):
    ```sh
    go get -d ./...
    ```

4. Place your static files (like images) in the `frontend/static/homePageImg` directory.

## Running the Application

Make sure you have the api running: 
1. Navigate to the `api/` directory:
    ```sh
    cd api
    ```

2. Start the api:
    ```sh
    make run
    ```

Now to start the application frontend:

1. Navigate to the `frontend/src` directory:
    ```sh
    cd frontend/src
    ```

2. Start the server:
    ```sh
    go run main.go
    ```

3. Open your web browser and go to `http://localhost:8080` to view the application.

## Endpoints

- **Home Page**: `GET /`
    - Displays a welcome message and a link to view cars.

- **View Cars**: `GET /cars`
    - Displays a list of cars.
    - Supports search by name and filter by manufacturer.

- **Car Details**: `GET /cars/{id}`
    - Displays detailed information about a specific car, including specifications and manufacturer details.


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.# cars
