package models

import (
	"fmt"
	"time"
)

type Reservation struct {
	ID         string    `json:"id" db:"id"`
	Email      string    `json:"email" db:"email"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	NationalId string    `json:"national_id" db:"national_id"`
	StartDate  time.Time `json:"start_date" db:"start_date"`
	EndDate    time.Time `json:"end_date" db:"end_date"`
	NumGuests  int       `json:"num_guests" db:"num_guests"`
}

func (r Reservation) String() string {
	return fmt.Sprintf(
		"Reservation(email: %s, firstName: %s, lastName: %s)",
		r.Email,
		r.FirstName,
		r.LastName,
	)
}