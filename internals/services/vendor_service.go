package services

import (
	DB "Medistock_Backend/internals/db"
	models "Medistock_Backend/internals/models"
	"fmt"
	"log"
	"strconv"
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

func RetrieveVendorByEmail(vendorEmail string) (*models.Vendor, error) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(RetrieveVendor)")
		return nil, fmt.Errorf("db instance is null.(RetrieveAllVendors)")
	}

	var QUERY string = `
		SELECT * from vendors where email = ?;
	`

	result, err := dbInstance.Query(QUERY, vendorEmail)
	if err != nil {
		log.Printf("Failed to retrieve vendor by email : %v", err)
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

// fetch vendor by its id
func RetrieveVendor(vendorId int, vendorEmail string) (*models.Vendor, error) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(RetrieveVendor)")
		return nil, fmt.Errorf("db instance is null.(RetrieveAllVendors)")
	}

	var QUERY string = `
		SELECT * FROM medistock_db.vendors where email=? OR id = ?;
	`

	result, err := dbInstance.Query(QUERY, vendorEmail,vendorId)
	if err != nil {
		log.Printf("Failed to retrieve vendor : %v", err)
		return nil, err
	}

	defer result.Close()
	var vendor models.Vendor

	for result.Next() {
		err := result.Scan(
			&vendor.ID,
			&vendor.Name,
			&vendor.ContactPerson,
			&vendor.Phone,
			&vendor.Email,
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
			&vendor.Phone,
			&vendor.Email,
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
	
	log.Println("[BACKEND] response : ",len(vendorList));
	return vendorList, nil
}

// vendor + supply
func UpsertSupplyItemService(supplyModel models.Supply, vendorId int, supplyPrice string) error {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewSupplyItemService)")
		return fmt.Errorf("db instance is null.(AddNewSupplyItemService)")
	}

	var QUERY string = `
		INSERT INTO supplies (id, name, sku, unit_of_measure, category, is_vital) 
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			sku = VALUES(sku),
			unit_of_measure = VALUES(unit_of_measure),
			category = VALUES(category),
			is_vital = VALUES(is_vital)
	`

	result, err := dbInstance.Exec(QUERY,
		supplyModel.ID,
		supplyModel.Name,
		supplyModel.SKU,
		supplyModel.UnitOfMeasure,
		supplyModel.Category,
		supplyModel.IsVital)

	if err != nil {
		log.Println("failed to execute upsert supply item service !")
		log.Println(err.Error())
		return err
	}

	// now add the supply item in supply-vendor table
	vendorModel, _ := RetrieveVendor(vendorId, "")
	amount, _ := strconv.ParseFloat(supplyPrice, 64)
	log.Printf("Upserting vendor_id=%d, supply_id=%s", vendorId, supplyModel.ID)

	QUERY = `
		INSERT INTO vendor_supply_prices (vendor_id,supply_id,unit_price,quality_rating,avg_delivery_days)
		VALUES (?,?,?,?,?) 
		ON DUPLICATE KEY UPDATE 
			unit_price=VALUES(unit_price),
			quality_rating=VALUES(quality_rating),
			avg_delivery_days=VALUES(avg_delivery_days)
	`

	result, err = dbInstance.Exec(QUERY,
		vendorId,
		supplyModel.ID,
		amount,
		vendorModel.OverallQualityRating,
		vendorModel.AvgDeliveryTimeDays)

	if err != nil {
		log.Println("Failed to upsert supply-combo-vendor record!")
		log.Println(err.Error())
		return err
	}

	resultId, err := result.LastInsertId()
	if err != nil {
		log.Println("Failed to get Last inserted Id !")
		log.Println(err.Error())
		return err
	}

	log.Println("Successfully added the supply-combo-vendor record.", resultId)
	return nil
}

func RetrieveSupply(supplyId string) (*models.Supply, error) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return nil, fmt.Errorf("db instance is null.(RetrieveAllVendors)")
	}

	var QUERY string = `
		SELECT * from supplies where id = ?;
	`

	result, err := dbInstance.Query(QUERY, supplyId)
	if err != nil {
		log.Printf("Failed to insert vendor : %v", err)
		return nil, err
	}

	defer result.Close()
	var supply models.Supply

	for result.Next() {
		err := result.Scan(
			&supply.ID,
			&supply.Name,
			&supply.SKU,
			&supply.UnitOfMeasure,
			&supply.Category,
			&supply.IsVital,
			&supply.CreatedAt,
			&supply.UpdatedAt,
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

	return &supply, nil
}

func AddNewVendorservice(vendorModel models.Vendor) error {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
		return fmt.Errorf("db instance is null.(AddNewVendorservice)")
	}

	log.Println("avg-delivery-time : ", vendorModel.AvgDeliveryTimeDays);
	log.Println("quality-rating : ", vendorModel.OverallQualityRating);

	var QUERY string = `
		INSERT IGNORE INTO vendors (name,contact_person,phone,email,address,overall_quality_rating,avg_delivery_time_days)
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


func FetchSuppliesByVendorId(vendorId int) []models.Supply {
	dbInstance := DB.Get()
	

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
	}

	var QUERY string = `
		SELECT * FROM medistock_db.supplies WHERE id IN (
			SELECT supply_id from medistock_db.vendor_supply_prices where
			vendor_id = ?
		);
	`;

	log.Println("vendor id is ", vendorId);
	result, err := dbInstance.Query(QUERY, vendorId);

	if err != nil {
		log.Fatal("Something went wrong while querying supplies. ", err.Error());
	}

	defer result.Close()

	var supplies []models.Supply;
	for result.Next() {
		var supply models.Supply;

		err := result.Scan(
			&supply.ID,
			&supply.Name,
			&supply.SKU,
			&supply.UnitOfMeasure,
			&supply.Category,
			&supply.IsVital,
			&supply.CreatedAt,
			&supply.UpdatedAt,
		)

		if (err!= nil ) {
			log.Fatal("Something went wrong while reading rows ,on querying supplies from vendors.",err.Error());
		}

		supplies= append(supplies, supply);
	}

	if err := result.Err(); err != nil {
		log.Fatalf("Row iteration error : %v", err)
	}

	log.Println("total supplies : ", len(supplies))

	return supplies;
}
