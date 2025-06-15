package handlers

import (
	models "Medistock_Backend/internals/models"
	services "Medistock_Backend/internals/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type status map[string]interface{}
func AddVendorHandler(w http.ResponseWriter, r *http.Request) {

	if (r.Method != http.MethodPost) {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be POST")
		return;
	}

	var vendor models.Vendor;
	validationController := validator.New()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vendor)
	if err != nil{
		// services.SetErrorResponse(w, http.StatusInternalServerError, "failed to parse request body")
		log.Printf("Something went wrong while creating new vendor : %v", err);
		log.Fatal(err)//exit
		return;
	}

	// check validations
	validationErr := validationController.Struct(&vendor)
	if validationErr != nil {
		w.Write([]byte("validations verification failed on parsed request body"))
		log.Fatal(validationErr)
	}

	// created, updated_at is handled by mysql
	log.Print("Adding new vendor to db !");
	err = services.AddNewVendorservice(vendor)
	if err!=nil {
		log.Fatal(err);
		return;
	}

	log.Println("Successfully inserted the record !");
	
	services.SetSuccessResponse(w, http.StatusOK);
	response := status{
		"status":http.StatusOK,
		"message":"success",
	}

	_ = json.NewEncoder(w).Encode(response)
	

}