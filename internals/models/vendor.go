package models

import (
	"time"
)

type Vendor struct {
	ID                   int       `json:"id" db:"id"` // Using int for auto-increment ID
	Name                 string    `json:"name" db:"name"`
	ContactPerson        string    `json:"contact_person" db:"contact_person"`
	Phone                string    `json:"phone" db:"phone"`
	Email                string    `json:"email" db:"email"`
	Address              string    `json:"address" db:"address"`
	OverallQualityRating float64   `json:"overall_quality_rating" db:"overall_quality_rating"`
	AvgDeliveryTimeDays  float64   `json:"avg_delivery_time_days" db:"avg_delivery_time_days"`
	Score                float64   `json:"score" db:"score"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}
