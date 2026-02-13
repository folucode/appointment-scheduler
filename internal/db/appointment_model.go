package db

import "time"

type Appointment struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	Description        string    `json:"description"`
	ContactInformation Contact   `json:"contact_information"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
}

type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
