package routers

import (
	handlers "Medistock_Backend/internals/handlers"
	middleware "Medistock_Backend/internals/middleware"
	"net/http"

	"github.com/gorilla/mux"
)


func RegisterHospitalRoutes(apiRouter *mux.Router) {
	
	apiRouter.Use(func(next http.Handler) http.Handler {
		return middleware.RateLimitMiddleWare(next.ServeHTTP);
	});

	apiRouter.HandleFunc("/hospital-client", handlers.AddHospitalHandler).Methods("POST");
	apiRouter.HandleFunc("/hospital-client/{id}", handlers.RetrieveUniqueHospital).Methods("GET");
	apiRouter.HandleFunc("/hospital-client",handlers.RetrieveHospitalByEmailHandler).Methods("GET");
	apiRouter.HandleFunc("/hospital-supplies", handlers.RetrieveVendorsHandler).Methods("GET");
    apiRouter.HandleFunc("/hospital-bulk-order-supplies", handlers.UpdateHospitalInventoryHandler).Methods("POST")
	apiRouter.HandleFunc("/hospital-bulk-order-results", handlers.RetreiveHospitalInventory).Methods("GET")
}