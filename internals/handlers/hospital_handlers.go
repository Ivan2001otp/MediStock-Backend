package handlers

import (
	models "Medistock_Backend/internals/models"
	services "Medistock_Backend/internals/services"
	"encoding/json"
	"log"
	"net/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func RetrieveUniqueHospital(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be GET")
		return
	}

	params := mux.Vars(r)

	hospitalId, _ := params["id"]

	hospitalClient, err := services.RetrieveHospital(hospitalId)
	if err != nil {
		log.Println("Something went wrong while retrieving HospitalClient by id")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	response := models.Message{
		"data":    hospitalClient,
		"status":  http.StatusOK,
		"message": "success",
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Println("Could not parse response - ", err)
		return
	}
}

func AddHospitalHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be POST")
		return
	}

	var hospitalClient models.Hospital
	validationController := validator.New()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&hospitalClient)

	if err != nil {
		log.Panic("Something went wrong on parsing request body - ", err)
		return
	}

	// check validation.
	validationErr := validationController.Struct(&hospitalClient)
	if validationErr != nil {
		w.Write([]byte("validations verification failed on parsed body"))
		log.Panic("validations verification failed on parsed body - ", validationErr)
		return
	}

	log.Printf("Adding new hospital client")
	err = services.AddNewHospitalClient(hospitalClient)
	if err != nil {
		log.Panic(err)
		return
	}

	log.Println("Successfully inserted hospital client !")
	services.SetSuccessResponse(w, http.StatusOK)

	response := models.Message{
		"status":  http.StatusOK,
		"message": "success",
	}

	_ = json.NewEncoder(w).Encode(response)
}
