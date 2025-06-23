package models

import "time"

type User struct {
	ID       int
	Email    string
	Password string
	Actor    string //VENDOR or HOSPITAL
}

type AuthToken struct {
	ID            string
	Email         string
	Actor         string // VENDOR or HOSPITAL
	Refresh_Token string
	Expiry_Time   time.Time
	Created_At    time.Time
}