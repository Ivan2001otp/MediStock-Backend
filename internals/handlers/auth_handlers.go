package handlers

import (
	models "Medistock_Backend/internals/models"
	services "Medistock_Backend/internals/services"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LoginHanlder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var request models.User
	json.NewDecoder(r.Body).Decode(request)

	access_token, refresh_token, err, status_code := services.ProcessAndGenerateTokenService(request)

	if err != nil {
		http.Error(w, err.Error(), status_code);
		return;
	}

	response := models.Message{
		"status":http.StatusOK,
		"message":"success",
		"data" : models.Message{
			"access_token":access_token,
			"refresh_token":refresh_token,
		},
	}

	json.NewEncoder(w).Encode(response);

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	json.NewDecoder(r.Body).Decode(user)

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	err = services.StoreUserService(user.Email, user.Actor, string(hashedPwd))

	if err != nil {
		log.Println("Something went wrong while adding Registration of new user !")
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := models.Message{
		"status":  http.StatusCreated,
		"message": "success",
		"data":    "User Registered - " + user.Email,
	}

	json.NewEncoder(w).Encode(response)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	// in request body we will send, the "refresh-token".
	// this is made ccall from frontend in the interceptor.
	if r.Method != http.MethodPost {
		http.Error(w, "Supposed to be POST", http.StatusBadRequest);
		return;
	}

	var request struct {
		Refresh_Token string `json:"refresh_token"`
	}

	json.NewDecoder(r.Body).Decode(&request);
	new_access_token ,err, statusCode :=  services.RenewAccessTokenService(request.Refresh_Token)
	if err !=  nil {
		http.Error(w, err.Error(), statusCode);
		return;
	}

	response := models.Message{
		"status":http.StatusOK,
		"message":"success",
		"data" : models.Message{
			"access_token":new_access_token,
		},
	}

	json.NewEncoder(w).Encode(response)

}
