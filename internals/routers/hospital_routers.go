package routers

import (
	middleware "Medistock_Backend/internals/middleware"
	handlers "Medistock_Backend/internals/handlers"
	"net/http"
	"github.com/gorilla/mux"
)


func RegisterHospitalRoutes(apiRouter *mux.Router) {
	apiRouter.Use(func(next http.Handler) http.Handler {
		return middleware.RateLimitMiddleWare(next.ServeHTTP);
	});

	apiRouter.HandleFunc("/hospital-client", handlers.AddHospitalHandler).Methods("POST");
	apiRouter.HandleFunc("/hospital-client/{id}", handlers.RetrieveUniqueHospital).Methods("GET")
}