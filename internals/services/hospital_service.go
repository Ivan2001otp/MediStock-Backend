package services

import (
	DB "Medistock_Backend/internals/db"
	models "Medistock_Backend/internals/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func UpdateHospitalSupplies(ctx context.Context, hospitalId string, supplyId string, quantityRecieved float64, reorderThreshold *float64) error {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(RetrieveHospitalByEmail)")
		return fmt.Errorf("db instance is null.(RetrieveHospitalByEmail)")
	}

	// start the transaction for atomicity.
	transaction, err := dbInstance.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal("Failed to begin transaction : ", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Println("rolling back , because error occured and panick !")

			transaction.Rollback()
			panic(p) // Rethrow panic after rollback.
		} else if err != nil {
			log.Println("rolling back , because error occured !")
			transaction.Rollback() // rollback if any error occured.
		} else {
			err = transaction.Commit()
			if err != nil {
				log.Printf("Failed to commit transaction")
				log.Fatal(err)
			}
		}
	}()

	// check for available global stock and lock the row to avoid race condition
	var currentGlobalStock float64
	var unit_of_measure_str string

	SELECT_STOCK_QUERY := `SELECT unit_of_measure FROM medistock_db.supplies where id = ? FOR UPDATE;`
	err = transaction.QueryRowContext(ctx, SELECT_STOCK_QUERY, supplyId).Scan(&unit_of_measure_str)

	log.Println("unit-of-measure-str : ", unit_of_measure_str)
	currentGlobalStock, _ = strconv.ParseFloat(strings.TrimSpace(unit_of_measure_str), 64)
	log.Println("global stock ", currentGlobalStock)
	log.Println("current requested stock ", quantityRecieved)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("supply with ID %s not found in global supplies table", supplyId)
		}

		log.Fatalf("failed to get current global stock for supply %s: %w", supplyId, err)
		return err
	}

	if currentGlobalStock < quantityRecieved {
		return fmt.Errorf("insufficient global stock for supply %s. Available: %.2f, Requested: %.2f", supplyId, currentGlobalStock, quantityRecieved)
	}

	// step 2 : deduct quantity from global supplies
	UPDATE_GLOBAL_STOCK_QUERY := `
		UPDATE medistock_db.supplies
		SET unit_of_measure = ?,
		updated_at = NOW()
		WHERE id = ?;
	`

	remaining_stocks := fmt.Sprintf("%v", (currentGlobalStock - quantityRecieved))
	result, err := transaction.ExecContext(ctx, UPDATE_GLOBAL_STOCK_QUERY, remaining_stocks, supplyId)
	if err != nil {
		log.Fatalf("failed to deduct quantity from global supplies for supply %s: %w", supplyId, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("failed to get rows affected after deducting from global supplies: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when deducting from global supplies for supply %s (possible concurrent update or ID not found)", supplyId)
	}

	log.Printf("Deducted %.2f units from global supply %s. New stock: %.2f", quantityRecieved, supplyId, currentGlobalStock-quantityRecieved)

	UPSERT_INVENTORY_QUERY := `
		INSERT INTO inventory (
			hospital_id,
			supply_id,
			current_stock,
			reorder_threshold
		) VALUES (
			?, ?, ?,
			COALESCE(?, DEFAULT(reorder_threshold))
		)
		ON DUPLICATE KEY UPDATE
			current_stock = current_stock + VALUES(current_stock);
	`

	_, err = transaction.ExecContext(ctx, UPSERT_INVENTORY_QUERY, hospitalId, supplyId, quantityRecieved, reorderThreshold)
	if err != nil {
		return fmt.Errorf("failed to record received supplies into hospital inventory for hospital %d, supply %s: %w", hospitalId, supplyId, err)
	}

	log.Printf("Successfully recorded/updated inventory for hospital %d, supply %s. Quantity added: %.2f", hospitalId, supplyId, quantityRecieved)
	return nil
}

func RetrieveHospitalByEmail(email string) (*models.Hospital, error) {
	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(RetrieveHospitalByEmail)")
		return nil, fmt.Errorf("db instance is null.(RetrieveHospitalByEmail)")
	}

	var QUERY string = `
		SELECT * from hospitals where contact_email = ?;
	`

	result, err := dbInstance.Query(QUERY, email)

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

func RetrieveHospitalInventory(ctx context.Context, hospitalId string, ) ([] models.HospitalInventoryItem) {
	inventoryItems := []models.HospitalInventoryItem {}


	dbInstance := DB.Get()

	if dbInstance == nil {
		log.Fatal("Db Instance is null.(AddNewVendorservice)")
	}

	QUERY := `
		SELECT
			inv.supply_id,
			s.name AS supply_name,
			s.sku AS supply_sku,
			s.category AS supply_category,
			s.is_vital,
			inv.current_stock,
			inv.reorder_threshold,
			v.name AS vendor_name 
		FROM 
			medistock_db.inventory inv
		JOIN
			medistock_db.supplies s ON inv.supply_id = s.id
		JOIN 
			medistock_db.vendor_supply_prices as vsp on vsp.supply_id = s.id
		JOIN
			medistock_db.vendors v on v.id =  vsp.vendor_id
		
		WHERE 
			inv.hospital_id = ?;
	`

	rows, err := dbInstance.QueryContext(ctx, QUERY, hospitalId);
	if err != nil {
		log.Fatalf("failed to execute query for hospital inventory for hospital %s: %v", hospitalId, err);

	}
	defer rows.Close();

	for rows.Next() {
		var item models.HospitalInventoryItem
		err := rows.Scan(
			&item.SupplyID,
			&item.SupplyName,
			&item.SupplySKU,
			&item.SupplyCategory,
			&item.IsVital,
			&item.CurrentStock,
			&item.ReorderThreshold,
			&item.VendorName,
		)
		if err != nil {
			log.Fatalf("failed to scan row into HospitalInventoryItem for hospital %d: %w", hospitalId, err)
		}
		inventoryItems = append(inventoryItems, item)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("error during row iteration for hospital inventory for hospital %d: %w", hospitalId, err)
	}

	log.Printf("Fetched %d inventory items for hospital %s", len(inventoryItems), hospitalId)
	
	return inventoryItems;
}