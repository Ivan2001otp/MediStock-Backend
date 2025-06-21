package models

import "time"

type Hospital struct {
	ID string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Address string `json:"address" db:"address"`
	ContactEmail string `json:"contact_email" db:"contact_email"`
	ContactPhone string `json:"contact_phone" db:"contact_phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time	`json:"updated_at" db:"updated_at"`
}