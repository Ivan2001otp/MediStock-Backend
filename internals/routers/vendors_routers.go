package routers

import (
	middleware "Medistock_Backend/internals/middleware"
	handlers "Medistock_Backend/internals/handlers"
	services "Medistock_Backend/internals/services"
	"net/http"
	"github.com/gorilla/mux"
)


func RegisterVendorRoutes(apiRouter *mux.Router) {

	apiRouter.Use(func(next http.Handler) http.Handler {
		return middleware.RateLimitMiddleWare(next, 2.0, 8)
	})

	apiRouter.HandleFunc("/health", services.HealthCheckHandler).Methods("GET")

	apiRouter.HandleFunc("/vendors",handlers.AddVendorHandler).Methods("POST")
	apiRouter.HandleFunc("/vendors", handlers.RetrieveVendorsHandler).Methods("GET")
	apiRouter.HandleFunc("/vendors/{id}", handlers.RetrieveUniqueVendor).Methods("GET")

	apiRouter.HandleFunc("/vendors/{id}", handlers.UpdateVendorHandler).Methods("PUT")

	
}
