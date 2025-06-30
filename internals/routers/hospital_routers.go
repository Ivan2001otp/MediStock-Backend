package routers

import (
	handlers "Medistock_Backend/internals/handlers"
	middleware "Medistock_Backend/internals/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


func RegisterHospitalRoutes(apiRouter *mux.Router) {
	
	apiRouter.Use(func(next http.Handler) http.Handler {
		return middleware.RateLimitMiddleWare(next.ServeHTTP);
	});

	log.Println("hospital routers invoked")
	apiRouter.HandleFunc("/hospital-client", handlers.AddHospitalHandler).Methods("POST");
	apiRouter.HandleFunc("/hospital-client/{id}", handlers.RetrieveUniqueHospital).Methods("GET");
	apiRouter.HandleFunc("/hospital-client",handlers.RetrieveHospitalByEmailHandler).Methods("GET");
	apiRouter.HandleFunc("/hospital-supplies", handlers.RetrieveVendorsHandler).Methods("GET");

}