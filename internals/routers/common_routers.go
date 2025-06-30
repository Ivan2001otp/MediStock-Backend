package routers
import (
	handlers "Medistock_Backend/internals/handlers"
	middleware "Medistock_Backend/internals/middleware"
	"net/http"
	"github.com/gorilla/mux"
)

func RegisterCommonRouters(apiRouter *mux.Router) {

	apiRouter.Use(func(next http.Handler) http.Handler {
		return middleware.RateLimitMiddleWare(next.ServeHTTP)
	});
	apiRouter.HandleFunc("/logout",handlers.LogoutHandler);
	apiRouter.HandleFunc("/vendors-supply/{id}", handlers.RetrieveSuppliesOfVendor).Methods("GET");
	
}