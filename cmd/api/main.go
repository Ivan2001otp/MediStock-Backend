package main

import (
	db "Medistock_Backend/internals/db"
	middleware "Medistock_Backend/internals/middleware"
	handlers "Medistock_Backend/internals/handlers"

	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Server started !")
	err := godotenv.Load()
	if err != nil {
		log.Println("Something went wrong. Could not load environment variables.")
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbConn, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}


	defer func() {
		_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := dbConn.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection pool closed gracefully.")
		}
	}()

	mainRouter := mux.NewRouter();
	apiRouter := mainRouter.PathPrefix("/api").Subrouter()
	apiRouter.Use(func(next http.Handler) http.Handler {
		return middleware.RateLimitMiddleWare(next,1.0, 5)
	})



	apiRouter.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	// Vendor Management Endpoints
	apiRouter.HandleFunc("/vendors", handlers.CreateVendorHandler).Methods("POST")
	apiRouter.HandleFunc("/vendors", handlers.GetAllVendorsHandler).Methods("GET")
	apiRouter.HandleFunc("/vendors/{id}", handlers.GetVendorByIDHandler).Methods("GET")
	apiRouter.HandleFunc("/vendors/{id}", handlers.UpdateVendorHandler).Methods("PUT")


	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("MediStock Go Backend API starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mainRouter))
}
