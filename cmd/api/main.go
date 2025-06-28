package main

import (
	db "Medistock_Backend/internals/db"
	"Medistock_Backend/internals/handlers"
	"Medistock_Backend/internals/middleware"
	routers "Medistock_Backend/internals/routers"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	log.Println("Server started !")
	err := godotenv.Load()
	if err != nil {
		log.Println("Something went wrong. Could not load environment variables.")
		return
	}

	// connecting to db
	rootCtx := context.Background()
	ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
	defer cancel()

	if err := db.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize DB : %v", err)
	}
	defer db.Close()

	// define cors config
	corsOptions := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedOrigins:   []string{"http://localhost:5173"},
	})

	// initializing all routers here
	mainRouter := mux.NewRouter()

	// public routes
	mainRouter.HandleFunc("/api/v1/login", handlers.LoginHandler).Methods("POST")
	mainRouter.HandleFunc("/api/v1/register", handlers.RegisterHandler).Methods("POST")
	mainRouter.HandleFunc("/api/v1/refresh-token", handlers.RefreshHandler).Methods("POST")

	// Adding routes for login,register
	vendorRouters := mainRouter.PathPrefix("/api/v1").Subrouter()
	hospitalRouters := mainRouter.PathPrefix("/api/v1").Subrouter()
	commonRouters := mainRouter.PathPrefix("/api/v1").Subrouter()

	// Apply Middlewares
	vendorRouters.Use(middleware.AuthMiddleware)
	vendorRouters.Use(middleware.CheckRoleMiddleware("VENDOR"))
	routers.RegisterVendorRoutes(vendorRouters)

	hospitalRouters.Use(middleware.AuthMiddleware)
	hospitalRouters.Use(middleware.CheckRoleMiddleware("HOSPITAL"))
	routers.RegisterHospitalRoutes(hospitalRouters)

	routers.RegisterCommonRouters(commonRouters)

	// setting handler with cors config.
	handler := corsOptions.Handler(mainRouter) // ?

	// setting our backend all prep !
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("MediStock Go Backend API starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
