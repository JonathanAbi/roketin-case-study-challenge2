# Challenge 2: Movie Festival API

This is the solution for Challenge 2 of the coding case study, a REST API for managing movie data for a short film festival.

## Technology Stack

* **Go (Golang)**: Primary programming language.
* **Chi**: As the HTTP router.
* **GORM**: As the ORM for database interaction.
* **MySQL**: As the relational database.

## Architecture

This project is built with an approach inspired by Clean Architecture, with responsibilities separated into the main layers:
* **Entity**: Core data structures.
* **Repository**: Abstraction for data access to the database.
* **Flow**: Application's business logic.
* **Parser**: Validates and transforms request data.
* **Handler**: Manages HTTP requests and responses.

## API Features

* **Create & Upload Movie**: `POST /api/movies`
    * Accepts movie metadata and video file via `multipart/form-data`.
* **Update Movie**: `PUT /api/movies/{id}`
    * Updates movie metadata via `application/x-www-form-urlencoded`.
* **List All Movies**: `GET /api/movies`
    * Supports pagination (`?page=...&limit=...`).
* **Search Movies**: `GET /api/movies/search`
    * Searches by title, description, genres, and artists. Supports pagination.
* **Delete Movie**: `DELETE /api/movies/{id}`
    * Uses soft delete.

## Setup and Running Instructions

1.  **Prerequisites:**
    * Go.
    * An active MySQL server.
    * Create an empty database in MySQL. Tables will be auto-migrated by GORM.

2.  **Configuration:**
    * Clone this repository:
        ```bash
        git clone https://github.com/JonathanAbi/roketin-case-study-challenge2
        cd roketin-case-study-challenge2
        ```
    * Create a `.env` file in the project's root directory.
    * Fill the `.env` file with your database configuration:
        ```env
        MYSQL_DSN="user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
        APP_PORT="8080"
        ```
        Replace `user`, `password`, `host`, `port`, and `dbname` with your MySQL setup details. `APP_PORT` is optional (defaults to 8080).

3.  **Running the Application:**
    * Ensure your Go module name is correctly referenced in all import paths. If you initialized with `go mod init [your_module_name]`, adjust import paths in the code accordingly.
    * Download dependencies:
        ```bash
        go mod tidy
        ```
    * Run the API server:
        ```bash
        go run main.go
        ```
    * The server will be running at `http://localhost:[APP_PORT]`.

## API Endpoint Summary

* `POST /api/movies`: Create a new movie (use `multipart/form-data` with fields `title`, `description`, `duration_minutes`, `artists`, `genres`, and `movieFile`).
* `GET /api/movies`: List all movies (use query params like `?page=1&limit=10`).
* `GET /api/movies/search`: Search movies (use query params like `?title=...&description=...&genre=...&artist=...&page=1&limit=10`).
* `PUT /api/movies/{id}`: Update a movie (send data as `application/x-www-form-urlencoded` or `multipart/form-data` if not updating file).
* `DELETE /api/movies/{id}`: Delete a movie.

---
