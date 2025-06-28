package main

import (
	DB "Medistock_Backend/internals/db"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ML_DATA struct {
	ID     int     `json:"id" db:"id"`
	email  string  `json:"email" db:"email"`
	rating float64 `json:"overall_quality_rating" db:"overall_quality_rating"`
	price  float64 `json:"unit_price" db:"unit_price"`
	days   float64 `json:"avg_delivery_days" db:"avg_delivery_days"`
}

func getAggregatedVendorDataBatch(ctx context.Context, limit int, lastProcessedId int) []ML_DATA {
	aggregatedDate := []ML_DATA{}

	query := `
		SELECT vendor.id , vendor.email , vendor.overall_quality_rating as rating,
		AVG(vendor_supply.unit_price) as price, AVG(vendor_supply.avg_delivery_days) as days
		FROM vendors vendor JOIN vendor_supply_prices vendor_supply
		ON vendor.id = vendor_supply.vendor_id
		WHERE vendor.id > ?
		GROUP BY vendor.id 
		ORDER By vendor.id 
		LIMIT ?;
	`

	dbInstance := DB.Get()
	if dbInstance == nil {
		log.Fatal("(CRON) Database instance not instantiated")
	}

	rows, err := dbInstance.Query(query, lastProcessedId, limit)
	if err != nil {
		log.Println("Could not execute GET query !")
		log.Fatal(err.Error())
	}

	defer rows.Close()
	for rows.Next() {
		var item ML_DATA

		err := rows.Scan(
			&item.ID,
			&item.email,
			&item.rating,
			&item.price,
			&item.days,
		)

		if err != nil {
			log.Fatalf("Error scanning target rows : %v", err)
		}

		aggregatedDate = append(aggregatedDate, item)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Something went wrong after row-scan : %v", err)
	}

	return aggregatedDate
}

func getMLWorthScore(ctx context.Context, batchData ML_DATA) (float64, error) {
	const ml_url string = "http://localhost:8000/predict"

	payload := map[string]interface{}{
		"unit_price":        batchData.price,
		"quality_rating":    int(batchData.rating),
		"avg_delivery_days": batchData.days,
	}

	jsonData, err := json.Marshal(payload)

	if err != nil {
		log.Fatal("Error marshalling json : %v", err)
	}

	response, err := http.Post(ml_url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatalf("Error on making POST request : %v", err)
	}

	defer response.Body.Close()

	/*	Response schema
		{
			"predicted_outcome_score": 0.5297
		}
	*/

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("Error reading response body : %v", err)
	}

	log.Println("Response status = ", response.Status)

	var responseJson map[string]interface{}

	if err = json.Unmarshal(body, &responseJson); err != nil {
		log.Println("Failed to unmarshal response-json : %v", err)
	}
	log.Println("response Json : ", responseJson)

	score_value, ok := responseJson["predicted_outcome_score"].(float64)
	if !ok {
		log.Println("Key not found or not a float64")
		return 0.0, fmt.Errorf("Invalid response format")
	}

	return score_value, nil
}

func updateVendorScore(ctx context.Context, batchId int, score float64) error {
	dbInstance := DB.Get()
	if dbInstance == nil {
		log.Fatal("(CRON) Database instance not instantiated")
	}

	query := `
		UPDATE medistock_db.vendors SET score = ? where id = ?;
	`

	_, err := dbInstance.ExecContext(ctx, query, score, batchId)
	if err != nil {
		log.Printf("Error updating vendor score (ID: %d): %v", batchId, err)
		return err
	}

	log.Printf("Successfully updated vendor (ID: %d) with score: %d", batchId, score)
	return nil
}

func runUpdateScores() {
	log.Println("Starting vendor ML worth score update process (Keyset Pagination)......")
	ctx := context.Background()

	lastProcessedVendorId := 0

	for {
		vendorDataBatch := getAggregatedVendorDataBatch(ctx, 5, lastProcessedVendorId)

		if len(vendorDataBatch) == 0 {
			log.Println("(CRON) - No more rows found !")
			break
		}

		for _, batchData := range vendorDataBatch {
			score, err := getMLWorthScore(ctx, batchData)

			if err != nil {
				log.Println("Error getting ML score. ", err.Error())
				continue
			}
			log.Println("Score from ML : ", score)

			err = updateVendorScore(ctx, batchData.ID, score)

			if err != nil {
				log.Println("Error updating score for vendor - ", err.Error())
			} else {
				log.Println("Successfully updated score for ", batchData.email)
			}

			lastProcessedVendorId = vendorDataBatch[len(vendorDataBatch)-1].ID
			time.Sleep(1 * time.Second)
		}

		log.Println("Vendor ML data score assignment process is completed.")
	}
}

func main() {
	log.Println("Starting cron...")
	/*
		- one by one filter the vendor by its id in supply-vendor-table
		- avg the price
		- send the the one tuple to ML model to give score for the given vendor-id.
		- this same process, goes to other vendors as well.
		- after every update have a delay...
	*/

	rootCtx := context.Background()
	ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		fmt.Println("Shutting down gracefully")
		cancel()
	}()

	if err := DB.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize DB : %v", err)
	}
	defer DB.Close()

	log.Println("(CRON) Waiting for cron to run next Batch .")
	time.Sleep(waitUntilNext3Hour())

	for {
		select {
		case <-ctx.Done():
			// Gracefull exit when signal is recieved.
			return

		default:
			runUpdateScores()
			select {
			case <-ctx.Done():
				return

			case <-time.After(3 * time.Hour):
				// wait before running task.
			}
		}
	}
}

func waitUntilNext3Hour() time.Duration {
	now := time.Now()

	next := now.Truncate(time.Hour).Add(time.Hour * time.Duration(3-(now.Hour()%3)))

	if next.Before(now) {
		next = next.Add(3 * time.Hour)
	}

	return time.Until(next)
}
