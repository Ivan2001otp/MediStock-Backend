package services

import (
	DB "Medistock_Backend/internals/db"
	models "Medistock_Backend/internals/models"
	"fmt"
	"log"
)

func RetrieveHospital(hospitalId string) (*models.Hospital, error) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return nil, fmt.Errorf("db instance is null.(RetrieveAllVendors)")
	}

	var QUERY string = `
		SELECT * from hospitals where id = ?;
	`
	result, err := dbInstance.Query(QUERY, hospitalId)

	if err != nil {
		log.Printf("Failed to insert vendor : %v", err)
		return nil, err
	}

	defer result.Close()
	var hospital models.Hospital

	for result.Next() {
		err := result.Scan(
			&hospital.ID,
			&hospital.Name,
			&hospital.Address,
			&hospital.ContactEmail,
			&hospital.ContactPhone,
			&hospital.CreatedAt,
			&hospital.UpdatedAt,
		)

		if err != nil {
			log.Printf("Error scanning hospital row : %v", err)
			return nil, err
		}
		break
	}

	if err := result.Err(); err != nil {
		log.Printf("Row iteration error : %v", err)
		return nil, err
	}

	return &hospital, nil
}

func AddNewHospitalClient(hospitalClient models.Hospital) error {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return fmt.Errorf("db instance is null.(AddNewVendorservice)")
	}

	var QUERY string = `
		INSERT INTO hospitals (id,name,address,contact_email,contact_phone) 
		VALUES (?,?,?,?,?)
	`

	result, err := dbInstance.Exec(QUERY, GenerateUUID(), hospitalClient.Name, hospitalClient.Address, hospitalClient.ContactEmail, hospitalClient.ContactPhone)

	if err != nil {
		log.Printf("Failed to insert hospital : %v", err)
		return err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		log.Printf("Failed to get last insert ID : %v", err)
		return err
	}

	log.Printf("Inserted hospitalclient with ID : %d", lastId)
	return nil

}
