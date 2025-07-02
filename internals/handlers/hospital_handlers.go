package handlers

import (
	models "Medistock_Backend/internals/models"
	services "Medistock_Backend/internals/services"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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

// needs a "email" as member of params
func RetrieveHospitalByEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be GET")
		return
	}

	params := r.URL.Query()
	hospitalEmail := params.Get("email")

	var hospital *models.Hospital

	hospital, err := services.RetrieveHospitalByEmail(hospitalEmail)
	if err != nil {
		log.Fatalf("Something went wrong on retrieving hospital by email : %v", err)
		return
	}

	response := models.Message{
		"status":  http.StatusOK,
		"message": "success",
		"data":    hospital,
	}

	_ = json.NewEncoder(w).Encode(response)

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

func UpdateHospitalInventoryHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be POST")
		return
	}

	var requestPaylaod map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestPaylaod)
	if err != nil {
		log.Panic("Something went wrong on parsing request body - ", err)
		return
	}

	log.Println("Request payload ")
	hospitalId := requestPaylaod["client_id"].(string)
	orderedSupplyList := requestPaylaod["supplies"]
	// vendor := requestPaylaod["vendor"];

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	var default_value float64 = 10
	var ordered_units_float_value float64 = 0
	var isSafe bool = false

	for _, supplyItem := range orderedSupplyList.([]interface{}) {
		supplyMap, ok := supplyItem.(map[string]interface{})

		if !ok {
			log.Println("One of the supplies is not valid map")
			continue
		}

		id, _ := supplyMap["id"].(string)
		value, _ := supplyMap["unit_of_measure"] // string datatype
		

		orderedUnits, ok := value.(string)
		if ok {
			ordered_units_float_value, _ = strconv.ParseFloat(orderedUnits, 64)
			isSafe = true
		} else {
			isSafe = false
		}

		if isSafe {
			err := services.UpdateHospitalSupplies(ctx, hospitalId, id, ordered_units_float_value, &default_value)

			if err != nil {
				http.Error(w, err.Error(), http.StatusAccepted)
				return
			}

			time.Sleep(1 * time.Second)
		}

	}

	w.WriteHeader(http.StatusOK)

	response := models.Message{
		"status":  http.StatusOK,
		"message": "success",
	}

	_ = json.NewEncoder(w).Encode(response)

}

func RetreiveHospitalInventory(w http.ResponseWriter, r *http.Request) {

	if (r.Method != http.MethodGet) {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be GET")
		return
	}

	hospitalId := r.URL.Query().Get("id")
	
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	
	inventoryItems := services.RetrieveHospitalInventory(ctx, hospitalId);

	w.WriteHeader(http.StatusOK);

	_ = json.NewEncoder(w).Encode(inventoryItems);
}

func RetrieveVendorsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be GET")
		return
	}

	lastSeenId, _ := strconv.Atoi(r.URL.Query().Get("after"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	log.Println("The last seen Id (after cursor) : ", lastSeenId)
	log.Println("The pagesize : ", pageSize)

	vendorsList, err := services.RetrieveAllVendors(lastSeenId, pageSize)

	if err != nil {
		log.Println("Unable to retrieve vendors : %v", err)
		http.Error(w, "Unable to retrieve vendors.", http.StatusInternalServerError)
		return
	}

	var nextCursor *int

	if len(vendorsList) >= 1 {
		val := vendorsList[len(vendorsList)-1].ID + 1
		nextCursor = &val
	} else {
		nextCursor = nil // optional, Go defaults to nil for pointers
	}

	paginatedResponse := status{
		"data":        vendorsList,
		"status":      http.StatusOK,
		"next_cursor": nextCursor,
	}

	err = json.NewEncoder(w).Encode(paginatedResponse)
	if err != nil {
		log.Println("Could not parse response - %v", err)
		return
	}
}
