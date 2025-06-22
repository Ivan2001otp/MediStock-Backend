package models

import (
	"time"
)

type VendorComboSupply struct {
	vendor_id int `json:"vendor_id" db:"vendor_id"`
	supply_id string `json:"supply_id" db:"supply_id"`
	unit_price float64 `json:"unit_price" db:"unit_price"`
	quality_rating float64 `json:"quality_rating" db:"quality_rating"`
	avg_delivery_days float64 `json:"avg_delivery_days" db:"avg_delivery_days"`
	created_at time.Time `json:"created_at" db:"created_at"`
	updated_at time.Time 	`json:"updated_at" db:"updated_at"`
}