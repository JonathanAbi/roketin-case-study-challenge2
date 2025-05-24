package main

import (
	"fmt"
	"log"
	"net/http"
	"roketin-case-study-challenge2/config"
	"roketin-case-study-challenge2/internal"
	"roketin-case-study-challenge2/internal/database"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	_, err = database.InitMySQLDB(cfg.GetDBDSN())
	if err != nil {
		log.Fatalf("Failed to initialize MySQL database: %v", err)
	}

	fmt.Println("MySQL database initialized successfully")
	
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string {
			"status": "UP",
			"message": "API is running",
		}

		internal.RespondWithJSON(w, http.StatusOK, payload)
	})

	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	fmt.Printf("Server running at http://localhost%s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}