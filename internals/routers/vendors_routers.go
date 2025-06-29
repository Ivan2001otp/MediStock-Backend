package routers

import (
	handlers "Medistock_Backend/internals/handlers"
	middleware "Medistock_Backend/internals/middleware"
	services "Medistock_Backend/internals/services"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterVendorRoutes(apiRouter *mux.Router) {

	apiRouter.Use(func(next http.Handler) http.Handler {
		return middleware.RateLimitMiddleWare(next.ServeHTTP)
	})

	apiRouter.HandleFunc("/health", services.HealthCheckHandler).Methods("GET")

	apiRouter.HandleFunc("/vendors", handlers.AddVendorHandler).Methods("POST")
	apiRouter.HandleFunc("/vendors", handlers.RetrieveVendorsHandler).Methods("GET")
	apiRouter.HandleFunc("/vendors/{id}", handlers.RetrieveUniqueVendor).Methods("GET")
	apiRouter.HandleFunc("/vendors/{id}", handlers.UpdateVendorHandler).Methods("PUT")

	// This endpoint helps to add supply from vendor.
	apiRouter.HandleFunc("/vendors/{id}", handlers.AddNewSupplyHandler).Methods("POST")

	// This endpoint helps to update supply details from vendor
	apiRouter.HandleFunc("/vendors/{id}", handlers.UpdateSupplyHandler).Methods("PATCH")

	apiRouter.HandleFunc("/vendors-supply/{id}", handlers.RetrieveSuppliesOfVendor).Methods("GET");

}
