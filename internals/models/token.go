package models

import "time"

type User struct {
	ID       int	`json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string	`json:"password" db:"password"`
	Actor    string `json:"actor" db:"actor"`//VENDOR or HOSPITAL
}

type AuthToken struct {
	ID            string    `json:"id" db:"id"`
	Email         string	`json:"email" db:"email"`
	Actor         string  	`json:"actor" db:"actor"`
	Refresh_Token string	`json:"refresh_token" db:"refresh_token"`
	Expiry_Time   time.Time	`json:"expiry_time" db:"expiry_time"`
	Created_At    time.Time	`json:"created_at" db:"created_at"`
}