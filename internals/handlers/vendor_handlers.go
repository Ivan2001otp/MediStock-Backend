package handlers

import (
	models "Medistock_Backend/internals/models"
	services "Medistock_Backend/internals/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type status map[string]interface{}

func UpdateVendorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be GET");
		return;
	}

	validationController := validator.New()
	params := mux.Vars(r);
	vendorId,_ := strconv.Atoi(params["id"])

	var updatedVendor models.Vendor;
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&updatedVendor)
	if err != nil {
		// services.SetErrorResponse(w, http.StatusInternalServerError, "failed to parse request body")
		log.Printf("Something went wrong while creating new vendor : %v", err)
		http.Error(w, "failed to parse request body", http.StatusInternalServerError);
		return
	}

	// check validations
	validationErr := validationController.Struct(&updatedVendor)
	if validationErr != nil {
		w.Write([]byte("validations verification failed on parsed request body"))
		http.Error(w, "validations on field failed gracefully", http.StatusInternalServerError);
		return;
	}


	oldVendor, err := services.RetrieveVendor(vendorId)
	if err != nil {
		log.Println("The vendor does not exists with vendorid-",vendorId);
		http.Error(w, "vendor does not exist", http.StatusInternalServerError);
		return;
	}

	
	// updating  fields
	oldVendor.Name = updatedVendor.Name
	oldVendor.ContactPerson = updatedVendor.ContactPerson
	oldVendor.Phone= updatedVendor.Phone
	oldVendor.Email = updatedVendor.Email
	oldVendor.Address = updatedVendor.Address
	oldVendor.OverallQualityRating = updatedVendor.OverallQualityRating
	oldVendor.AvgDeliveryTimeDays = updatedVendor.AvgDeliveryTimeDays

	err = services.UpdateVendor(*oldVendor);
	if err != nil {
		http.Error(w, "Upsert failed ", http.StatusInternalServerError);
		return;
	}

	response := status{
		"data":"success",
		"status":http.StatusOK,
	}

	_ = json.NewEncoder(w).Encode(response);
}

// get vendor by id.
func RetrieveUniqueVendor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be GET")
		return
	}

	params := mux.Vars(r)
	vendorId, _ := strconv.Atoi(params["id"])

	vendorModel, err := services.RetrieveVendor(vendorId)
	if err != nil {
		log.Println("Something went wrong while retrieving vendor by id")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	response := status{
		"data":   vendorModel,
		"status": http.StatusOK,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Could not parse response - %v", err)
		return
	}
}

func RetrieveVendorsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be GET")
		return
	}

	lastSeenId, _ := strconv.Atoi(r.URL.Query().Get("after"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	vendorsList, err := services.RetrieveAllVendors(lastSeenId, pageSize)

	if err != nil {
		log.Println("Unable to retrieve vendors : %v", err)
		http.Error(w, "Unable to retrieve vendors.", http.StatusInternalServerError)
		return
	}

	paginatedResponse := status{
		"data":        vendorsList,
		"status":      http.StatusOK,
		"next_cursor": vendorsList[len(vendorsList)-1].ID + 1,
	}

	err = json.NewEncoder(w).Encode(paginatedResponse)
	if err != nil {
		log.Println("Could not parse response - %v", err)
		return
	}
}

func AddVendorHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be POST")
		return
	}

	var vendor models.Vendor
	validationController := validator.New()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vendor)
	if err != nil {
		// services.SetErrorResponse(w, http.StatusInternalServerError, "failed to parse request body")
		log.Printf("Something went wrong while creating new vendor : %v", err)
		// log.Error(err) //exit
		return
	}

	// check validations
	validationErr := validationController.Struct(&vendor)
	if validationErr != nil {
		w.Write([]byte("validations verification failed on parsed request body"))
		log.Panic(validationErr)
		return;
	}

	// created, updated_at is handled by mysql
	log.Print("Adding new vendor to db !")
	err = services.AddNewVendorservice(vendor)
	if err != nil {
		log.Panic(err)
		return
	}

	log.Println("Successfully inserted the record !")

	services.SetSuccessResponse(w, http.StatusOK)
	response := status{
		"status":  http.StatusOK,
		"data": "success",
	}

	_ = json.NewEncoder(w).Encode(response)
}
