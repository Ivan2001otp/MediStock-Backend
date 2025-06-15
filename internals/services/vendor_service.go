package services

import (
	DB "Medistock_Backend/internals/db"
	models "Medistock_Backend/internals/models"
	"fmt"
	"log"
)

func AddNewVendorservice(vendorModel models.Vendor) (error){
	dbInstance := DB.Get();

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)");
		return  fmt.Errorf("db instance is null.(AddNewVendorservice)");
	}

	var QUERY string = `
		INSERT INTO vendors (name,contact_person,phone,email,address,overall_quality_rating,avg_delivery_time)
		 VALUES (?,?,?,?,?,?,?)
		`

	result, err := dbInstance.Exec(QUERY, vendorModel.Name, vendorModel.ContactPerson, vendorModel.Phone, vendorModel.Email, vendorModel.Address,vendorModel.OverallQualityRating, vendorModel.AvgDeliveryTimeDays);

	if err != nil {
		log.Printf("Failed to insert vendor : %v",err);
		return err;
	}

	lastId , err := result.LastInsertId();
	if err != nil {
		log.Printf("Failed to get last insert Id : %v", err)
		return err;
	}

	log.Printf("Inserted vendor with ID : %d", lastId)
	return nil;

}