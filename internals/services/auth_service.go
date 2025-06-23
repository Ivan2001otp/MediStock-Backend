package services

import (
	DB "Medistock_Backend/internals/db"
	models "Medistock_Backend/internals/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func RenewAccessTokenService(refreshToken string) (*string, error, int) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return nil, fmt.Errorf("db instance is null.(RetrieveAllVendors)"), http.StatusInternalServerError
	}

	var email, actor string
	var expiry time.Time

	err := dbInstance.QueryRow(`SELECT email,actor,expiry_time from auth_token where refresh_token = ?`, refreshToken).Scan(&email, &actor, &expiry)

	if err != nil {
		log.Println("Something went wrong on renewing fresh accesstokens !")
		return nil, err, http.StatusInternalServerError
	}

	if time.Now().After(expiry) {
		return nil, fmt.Errorf("Expired refresh token"), http.StatusUnauthorized
	}

	newAccessToken, err := GenerateAccessToken(email, actor)
	if err != nil {
		log.Println("Something went wrong while generating new access tokens (auth_service.go)")
		return nil, err, http.StatusInternalServerError
	}

	return &newAccessToken, nil, http.StatusOK
}

func ProcessAndGenerateTokenService(user models.User) (*string, *string, error, int) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return nil, nil, fmt.Errorf("db instance is null.(RetrieveAllVendors)"), http.StatusInternalServerError
	}

	var hashedPassword string
	err := dbInstance.QueryRow("SELECT password FROM users WHERE email = ? AND actor = ?", user.Email, user.Actor).Scan(&hashedPassword)

	if err != nil {
		log.Println("WARNING : Something went wrong while searching for exising user .")
		log.Println(err.Error())
		return nil, nil, fmt.Errorf("Invalid credentials"), http.StatusUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		return nil, nil, fmt.Errorf("Invalid Password"), http.StatusUnauthorized

	}

	access_token, err := GenerateAccessToken(user.Email, user.Actor)
	if err != nil {
		log.Println("Something went wrong ,when access_token was generated!")

		return nil, nil, err, http.StatusInternalServerError
	}
	refresh_token := GenerateRefreshToken()

	expiry := time.Now().Add(2 * 24 * time.Hour)
	_, err = dbInstance.Exec(`INSERT INTO auth_token (email,actor,refresh_token,expiry_time) VALUES (?,?,?,?)
		ON DUPLICATE KEY UPDATE 
		refresh_token = VALUES(refresh_token),
		expiry_time = VALUES(expiry_time)
	`, user.Email, user.Actor, refresh_token, expiry)

	if err != nil {
		log.Println("something went wrong while saving refresh-tokens.")
		return nil, nil, err, http.StatusInternalServerError
	}

	return &access_token, &refresh_token, nil, http.StatusOK

}

func StoreUserService(email, actor, hashedPass string) error {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return fmt.Errorf("db instance is null.(RetrieveAllVendors)")
	}

	var QUERY string = `INSERT INTO users (email,actor,password) values (?,?,?) 
	ON DUPLICATE KEY UPDATE
			email = VALUES(email),
			actor = VALUES(actor),
			password = VALUES(password) 
	`
	_, err := dbInstance.Exec(QUERY, email, actor, hashedPass)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println("Successfully added the new user in USERS table.")
	return nil
}

// access token generation
func GenerateAccessToken(email, actor string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"actor": actor,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("Generating claims and signing...")
	return token.SignedString([]byte(access_secret_key))
}

func GenerateRefreshToken() string {
	return GenerateUUID()
}

func GetSecretKey() string {
	return access_secret_key
}

var access_secret_key string = `b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAACFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAgEAuHtgQodovWKA2kCQn326nwJubT5yH5KcrHPStOYs9l9Crr0o9YUB`
