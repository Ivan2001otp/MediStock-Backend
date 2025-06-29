package services

import (
	DB "Medistock_Backend/internals/db"
	models "Medistock_Backend/internals/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func LogoutService(email, actor string) error {
	//actor means VENDOR or HOSPITAL.
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return fmt.Errorf("db instance is null.(RetrieveAllVendors)")
	}

	var QUERY string = `DELETE FROM auth_token where email=? and actor=?`
	_, err := dbInstance.Exec(QUERY, email, actor)
	if err != nil {
		log.Println("Db Error .Error during logout : ", err.Error())

		return err
	}

	return nil
}

func FetchUserByEmail(email string) (*models.User, error) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return nil, fmt.Errorf("db instance is null.(RetrieveAllVendors)")
	}

	var QUERY string = `SELECT id,actor,password from users where email = ?;`
	var id int
	var actor, hashPassword string = "", ""

	err := dbInstance.QueryRow(QUERY, email).Scan(&id, &actor, &hashPassword)
	if err != nil {
		log.Println("The record does not exists in DB for email : ", email)
		return nil, fmt.Errorf("Something happened while searching user by email")
	}

	if actor == "" && hashPassword == "" {
		return nil, fmt.Errorf("No user with %s exists", email)
	}

	var user models.User
	user.ID = id
	user.Actor = actor
	user.Email = email
	user.Password = hashPassword

	return &user, nil
}

func RenewAccessTokenService(refreshToken string) (*string, error, int) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return nil, fmt.Errorf("db instance is null.(RetrieveAllVendors)"), http.StatusInternalServerError
	}

	var email, actor string
	var expiry time.Time

	log.Println("refresh-token to be renewed : ",refreshToken);
	err := dbInstance.QueryRow(`SELECT email,actor,expiry_time from auth_token where refresh_token = ?`, refreshToken).Scan(&email, &actor, &expiry)

	if err != nil {
		log.Fatalf("Something went wrong on renewing fresh accesstokens : %v", err)
		return nil, err, http.StatusInternalServerError
	}

	if time.Now().After(expiry) {
		return nil, fmt.Errorf("Expired refresh token"), http.StatusForbidden // 403
	}

	newAccessToken, err := GenerateAccessToken(email, actor)
	if err != nil {
		log.Println("Something went wrong while generating new access tokens (auth_service.go)")
		return nil, err, http.StatusInternalServerError
	}

	log.Println("Generated new access token !");
	return &newAccessToken, nil, http.StatusOK
}

func ProcessAndGenerateTokenService(user models.User) (*string, *string, error, int) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return nil, nil, fmt.Errorf("db instance is null.(RetrieveAllVendors)"), http.StatusInternalServerError
	}

	access_token, err := GenerateAccessToken(user.Email, user.Actor)
	if err != nil {
		log.Println("Something went wrong ,when access_token was generated!");
		return nil, nil, err, http.StatusInternalServerError;
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
		log.Fatal(err)
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
	log.Println("Generated refresh token")
	return GenerateUUID()
}

func GetSecretKey() string {
	return access_secret_key
}

var access_secret_key string = `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`
