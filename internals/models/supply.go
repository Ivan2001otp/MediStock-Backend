package models

import (
	"time"
)

type Supply struct {
	ID string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	SKU string `json:"sku" db:"sku"`
	UnitOfMeasure string `json:"unit_of_measure" db:"unit_of_measure"`
	Category string `json:"category" db:"category"`
	IsVital bool `json:"is_vital" db:"is_vital"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}