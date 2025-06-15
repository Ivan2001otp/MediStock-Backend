package routers

import (
	middleware "Medistock_Backend/internals/middleware"
	services "Medistock_Backend/internals/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterVendorRoutes(apiRouter *mux.Router) {

	apiRouter.Use(func(next http.Handler) http.Handler {
		log.Println("request limit crossed !")
		return middleware.RateLimitMiddleWare(next, 1.0, 5)
	})

	apiRouter.HandleFunc("/health", services.HealthCheckHandler).Methods("GET")
	// apiRouter.HandleFunc("/vendors",handlers.AddnewVendorHandler).Methods("POST")
	// apiRouter.HandleFunc("/vendors", handlers.RetrieveallVendorsHandler).Methods("GET")
	// apiRouter.HandleFunc("/vendors/{id}", handlers.RetrieveVendorHandler).Methods("GET")

	// apiRouter.HandleFunc("/vendors/{id}", handlers.UpdateVendorHandler).Methods("PUT")

	
}
