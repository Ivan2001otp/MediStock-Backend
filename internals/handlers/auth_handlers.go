package handlers

import (
	models "Medistock_Backend/internals/models"
	services "Medistock_Backend/internals/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// "email","actor" as query params
	email := r.URL.Query().Get("email")
	actor := r.URL.Query().Get("actor") //VENDOR or HOSPITAL

	err := services.LogoutService(email, actor)
	if err != nil {
		log.Println("Something went wrong on logging out user !")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	message := models.Message{
		"status":  http.StatusOK,
		"message": "log out success",
	}

	_ = json.NewEncoder(w).Encode(message)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var request models.User
	var err error
	json.NewDecoder(r.Body).Decode(&request)

	log.Println("email : ", request.Email)
	log.Println("pass : ", request.Password)
	log.Println("role : ", request.Actor)

	fetchedUser, err := services.FetchUserByEmail(request.Email)
	if err != nil {
		log.Println("user for the given email does not exist")
		http.Error(w, err.Error(), http.StatusNoContent) //204 status
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(request.Password))

	if err != nil {
		log.Println("The password is not same. Thus its invalid !")
		http.Error(w, "Password is invalid", http.StatusNoContent)
		return
	}

	access_token, refresh_token, err, status_code := services.ProcessAndGenerateTokenService(request)
	if err != nil {
		http.Error(w, err.Error(), status_code)
		return
	}

	var response models.Message

	response = models.Message{
		"status":  http.StatusOK,
		"message": "success",
		"data": map[string]interface{}{
			"access_token":  access_token,
			"refresh_token": refresh_token,
			"actor":         fetchedUser.Actor,
			"email":         fetchedUser.Email,
		},
	}

	_ = json.NewEncoder(w).Encode(response)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	validationController := validator.New()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("email : ", user.Email)
	log.Println("password : ", user.Password)
	log.Println("actor : ", user.Actor)

	validationErr := validationController.Struct(&user)
	if validationErr != nil {
		w.Write([]byte("validations verification failed on parsed body"))
		log.Panic("validations verification failed on parsed body - ", validationErr)
		return
	}

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

	access_token, refresh_token, err, _ := services.ProcessAndGenerateTokenService(user)

	if err != nil {
		log.Println("Something went wrong while registration new user. Could not generate tokens !", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := models.Message{
		"status":  http.StatusOK,
		"message": "success",
		"data": map[string]interface{}{
			"access_token":  access_token,
			"refresh_token": refresh_token,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	// in request body we will send, the "refresh-token".
	// this is made ccall from frontend in the interceptor.
	if r.Method != http.MethodPost {
		http.Error(w, "Supposed to be POST", http.StatusBadRequest)
		return
	}

	var request struct {
		Refresh_Token string `json:"refresh_token"`
	}

	json.NewDecoder(r.Body).Decode(&request)
	new_access_token, err, statusCode := services.RenewAccessTokenService(request.Refresh_Token)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := models.Message{
		"status":  http.StatusOK,
		"message": "success",
		"data": models.Message{
			"access_token": new_access_token,
		},
	}

	json.NewEncoder(w).Encode(response)

}
