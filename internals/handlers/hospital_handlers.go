package handlers

import (
	models "Medistock_Backend/internals/models"
	
	"log"
	"net/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	services "Medistock_Backend/internals/services"
	
)

func AddHospitalHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be POST")
		return;
	}

	
}