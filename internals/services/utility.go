package services

import (
	"log"
	"net/http"
	"encoding/json"
)

type status map[string]interface{}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := status {
		"status" : 200,
		"message" : "I am alright !",
	}

	err := json.NewEncoder(w).Encode(response);

	if err != nil {
		SetErrorResponse(w, 500, "Health Checker is Down !")
		return;
	}

	SetSuccessResponse(w, 200);
	
}

func SetErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json");
	w.WriteHeader(statusCode);

	response := status {
		"status":statusCode,
		"error":message,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Unable to parse error response to json object !")
		http.Error(w, "unable to parse error response", http.StatusInternalServerError)
	}
}

func SetSuccessResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json");
	w.WriteHeader(statusCode);

}