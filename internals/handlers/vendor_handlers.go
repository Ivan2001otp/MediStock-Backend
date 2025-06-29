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

func UpdateSupplyHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPatch {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be PATCH")
		return
	}

	params := mux.Vars(r)
	vendorId, _ := strconv.Atoi(params["id"]) ///integer

	query_params := r.URL.Query()
	supplyId := query_params["supply_id"] //string

	updated_supply_price := query_params["supply_price"]
	updated_supply_name := query_params["supply_name"]
	updated_supply_sku := query_params["supply_sku"]

	updated_supply_unitofmeasure := query_params["supply_unit_of_measure"]
	updated_supply_category := query_params["supply_category"]
	updated_supply_isvital := query_params["supply_is_vital"]

	oldSupplyModel, err := services.RetrieveSupply(supplyId[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	oldSupplyModel.Name = updated_supply_name[0]
	oldSupplyModel.UnitOfMeasure = updated_supply_unitofmeasure[0]
	oldSupplyModel.SKU = updated_supply_sku[0]
	oldSupplyModel.Category = updated_supply_category[0]
	oldSupplyModel.IsVital, err = strconv.ParseBool(updated_supply_isvital[0])

	if err != nil {
		log.Panic("Failed to parse Bool (IsVital) supply model !")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = services.UpsertSupplyItemService(*oldSupplyModel, vendorId, updated_supply_price[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Successfully updated the supply - ", supplyId)
	response := models.Message{
		"message": "success",
		"status":  http.StatusOK,
	}

	_ = json.NewEncoder(w).Encode(response)
}

// Add new supply by particular vendor-id
func AddNewSupplyHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be GET")
		return
	}

	/*
		{
			"id":"",
			"name":"",
			"sku":"",
			"unit_of_measure":"",
			"category":"",
			"is_vital":"",
			"created_at",
			"updated_at"
		}
	*/

	// id is "path_param" , supply_price is "query_param"

	validationController := validator.New()
	params := mux.Vars(r)
	vendorId, _ := strconv.Atoi(params["id"])

	query_values := r.URL.Query()
	// calculated sum of supplyPrice will come from frontend itself.
	supplyPrice := query_values.Get("supply_price")

	var supplyModel models.Supply
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&supplyModel)
	if err != nil {
		log.Printf("Something went wrong while Adding new supply : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validationErr := validationController.Struct(&supplyModel)
	if validationErr != nil {
		w.Write([]byte("validations verification failed on parsed request body"))
		log.Panic(validationErr)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	supplyModel.ID = services.GenerateUUID()
	err = services.UpsertSupplyItemService(supplyModel, vendorId, supplyPrice)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("successfully inserted supply by vendorid ", vendorId)
	response := models.Message {
		"status" : http.StatusOK,
		"message" : "success",
	}

	_ = json.NewEncoder(w).Encode(response);
}

func UpdateVendorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be GET")
		return
	}

	validationController := validator.New()
	params := mux.Vars(r)
	vendorId, _ := strconv.Atoi(params["id"])

	var updatedVendor models.Vendor
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&updatedVendor)
	if err != nil {
		// services.SetErrorResponse(w, http.StatusInternalServerError, "failed to parse request body")
		log.Printf("Something went wrong while creating new vendor : %v", err)
		http.Error(w, "failed to parse request body", http.StatusInternalServerError)
		return
	}

	// check validations
	validationErr := validationController.Struct(&updatedVendor)
	if validationErr != nil {
		w.Write([]byte("validations verification failed on parsed request body"))
		http.Error(w, "validations on field failed gracefully", http.StatusInternalServerError)
		return
	}

	oldVendor, err := services.RetrieveVendor(vendorId, "")//email can be empty, but vendorId must not
	if err != nil {
		log.Println("The vendor does not exists with vendorid-", vendorId)
		http.Error(w, "vendor does not exist", http.StatusInternalServerError)
		return
	}

	// updating  fields
	oldVendor.Name = updatedVendor.Name
	oldVendor.ContactPerson = updatedVendor.ContactPerson
	oldVendor.Phone = updatedVendor.Phone
	oldVendor.Email = updatedVendor.Email
	oldVendor.Address = updatedVendor.Address
	oldVendor.OverallQualityRating = updatedVendor.OverallQualityRating
	oldVendor.AvgDeliveryTimeDays = updatedVendor.AvgDeliveryTimeDays

	err = services.UpdateVendor(*oldVendor)
	if err != nil {
		http.Error(w, "Upsert failed ", http.StatusInternalServerError)
		return
	}

	response := status{
		"data":   "success",
		"status": http.StatusOK,
	}

	_ = json.NewEncoder(w).Encode(response)
}

// get vendor by id.
func RetrieveUniqueVendor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SetErrorResponse(w, http.StatusBadRequest, "supposed to be GET")
		return
	}
	
	params := mux.Vars(r)
	vendorId, _ := strconv.Atoi(params["id"])


	// send email in query-params
	query_params := r.URL.Query()
	vendorEmail := query_params.Get("email")

	vendorModel, err := services.RetrieveVendor(vendorId, vendorEmail)
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
		log.Printf("Something went wrong while creating new vendor : %v", err)
		return
	}

	// check validations
	validationErr := validationController.Struct(&vendor)
	if validationErr != nil {
		w.Write([]byte("validations verification failed on parsed request body"))
		log.Panic(validationErr)
		return
	}

	// created, updated_at is handled by mysql
	log.Print("Adding new vendor to db !")
	log.Println("AVG - ", vendor.AvgDeliveryTimeDays);
	log.Println("RATING - ", vendor.OverallQualityRating);

	err = services.AddNewVendorservice(vendor)
	if err != nil {
		log.Panic(err)
		return
	}

	log.Println("Successfully inserted the record !")

	services.SetSuccessResponse(w, http.StatusOK)
	response := status{
		"status": http.StatusOK,
		"data":   "success",
	}

	_ = json.NewEncoder(w).Encode(response)
}

// Retrieve supplies of a corresponding vendors.
// The id of vendor is given to do the task.
func RetrieveSuppliesOfVendor(w http.ResponseWriter, r * http.Request) {

	if (r.Method != http.MethodGet) {
		http.Error(w, "supposed to be GET request", http.StatusBadRequest);
		return;
	}

	params := mux.Vars(r);
	vendorId,_ := strconv.Atoi(params["id"]);

	
	var supplies []models.Supply = services.FetchSuppliesByVendorId(vendorId);

	response := models.Message {
		"status":http.StatusOK,
		"message":"success",
		"data":supplies,
	}

	w.WriteHeader(http.StatusOK);
	_ = json.NewEncoder(w).Encode(response);

}
