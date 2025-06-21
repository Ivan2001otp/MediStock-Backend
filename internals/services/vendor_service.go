package services

import (
	DB "Medistock_Backend/internals/db"
	models "Medistock_Backend/internals/models"
	"fmt"
	"log"
)

func UpdateVendor(updatedVendor models.Vendor) error {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("DB instance is null . (UpdateVendor)")
		return fmt.Errorf("db instance is null.(UpdateVendor)")
	}

	query :=
		`
	INSERT INTO vendors(id,name,contact_person,phone,email,address,overall_quality_rating,avg_delivery_time_days)
	VALUES (?,?,?,?,?,?,?, ?)
	ON DUPLICATE KEY UPDATE
		id = VALUES(id),
		name = VALUES(name),
		contact_person = VALUES(contact_person),
		phone = VALUES(phone),
		email = VALUES(email),
		address = VALUES(address),
		overall_quality_rating = VALUES(overall_quality_rating),
		avg_delivery_time_days = VALUES(avg_delivery_time_days),
		updated_at = CURRENT_TIMESTAMP;
	`

	_, err := dbInstance.Exec(query, updatedVendor.ID, updatedVendor.Name, updatedVendor.ContactPerson, updatedVendor.Phone, updatedVendor.Email, updatedVendor.Address, updatedVendor.OverallQualityRating, updatedVendor.AvgDeliveryTimeDays)

	if err != nil {
		log.Printf("Upsert failed in UpdateVendor: %v", err)
		return err
	}

	return nil
}

// fetch vendor by its id
func RetrieveVendor(vendorId int) (*models.Vendor, error) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return nil, fmt.Errorf("db instance is null.(RetrieveAllVendors)")
	}

	var QUERY string = `
		SELECT * from vendors where id = ?;
	`

	result, err := dbInstance.Query(QUERY, vendorId)
	if err != nil {
		log.Printf("Failed to insert vendor : %v", err)
		return nil, err
	}

	defer result.Close()
	var vendor models.Vendor

	for result.Next() {
		err := result.Scan(
			&vendor.ID,
			&vendor.Name,
			&vendor.ContactPerson,
			&vendor.Email,
			&vendor.Phone,
			&vendor.Address,
			&vendor.OverallQualityRating,
			&vendor.AvgDeliveryTimeDays,
			&vendor.Score,
			&vendor.CreatedAt,
			&vendor.UpdatedAt,
		)

		if err != nil {
			log.Printf("Error scanning vendor row : %v", err)
			return nil, err
		}

		break
	}

	if err := result.Err(); err != nil {
		log.Printf("Row iteration error : %v", err)
		return nil, err
	}

	return &vendor, nil
}

// Fetch all vendors in paginated manner
func RetrieveAllVendors(lastSeenId int, pageSize int) ([]models.Vendor, error) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return nil, fmt.Errorf("db instance is null.(RetrieveAllVendors)")
	}

	if pageSize <= 0 || pageSize >= 50 {
		pageSize = 15 //safety default
	}

	var QUERY string = `
		SELECT id,name,contact_person,phone,email,address,overall_quality_rating,avg_delivery_time_days,score,created_at,updated_at FROM vendors
		WHERE id > ?
		ORDER By id ASC
		LIMIT ?
	`

	rows, err := dbInstance.Query(QUERY, lastSeenId, pageSize)
	if err != nil {
		log.Printf("Something went wrong while selecting all vendors !")
		return nil, err
	}

	defer rows.Close()
	var vendorList []models.Vendor

	for rows.Next() {
		var vendor models.Vendor

		err := rows.Scan(
			&vendor.ID,
			&vendor.Name,
			&vendor.ContactPerson,
			&vendor.Email,
			&vendor.Phone,
			&vendor.Address,
			&vendor.OverallQualityRating,
			&vendor.AvgDeliveryTimeDays,
			&vendor.Score,
			&vendor.CreatedAt,
			&vendor.UpdatedAt,
		)

		if err != nil {
			log.Printf("Error scanning vendor row : %v", err)
			return nil, err
		}

		vendorList = append(vendorList, vendor)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row iteration error : %v", err)
		return nil, err
	}

	return vendorList, nil

}

func AddNewVendorservice(vendorModel models.Vendor) error {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return fmt.Errorf("db instance is null.(AddNewVendorservice)")
	}

	var QUERY string = `
		INSERT INTO vendors (name,contact_person,phone,email,address,overall_quality_rating,avg_delivery_time_days)
		VALUES (?,?,?,?,?,?,?)
		`

	result, err := dbInstance.Exec(QUERY, vendorModel.Name, vendorModel.ContactPerson, vendorModel.Phone, vendorModel.Email, vendorModel.Address, vendorModel.OverallQualityRating, vendorModel.AvgDeliveryTimeDays)

	if err != nil {
		log.Printf("Failed to insert vendor : %v", err)
		return err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last insert Id : %v", err)
		return err
	}

	log.Printf("Inserted vendor with ID : %d", lastId)
	return nil

}
