package handlers

import (
	models "Medistock_Backend/internals/models"
	services "Medistock_Backend/internals/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func LoginHanlder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed);
		return;
	}

	var request models.User
	json.NewDecoder(r.Body).Decode(request);
	
	// Check whether the entered user exists in db or not.
	var user models.User
	access_token,refresh_token,err := services.ProcessAndGenerateTokens(request)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed);
		return;
	}

	var user models.User
	json.NewDecoder(r.Body).Decode(user);

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to hash password", http.StatusInternalServerError);
		return;
	}

	err = services.StoreUserService(user.Email, user.Actor, string(hashedPwd))

	if err != nil {
		log.Println("Something went wrong while adding Registration of new user !");
		log.Println(err.Error());
		http.Error(w, err.Error(), http.StatusInternalServerError);
		return;
	}

	w.WriteHeader(http.StatusCreated)
	response := models.Message{
		"status" : http.StatusCreated,
		"message":"success",
		"data":"User Registered - "+user.Email,
	}

	json.NewEncoder(w).Encode(response)
 }