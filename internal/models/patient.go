package models

import "time"

type Patient struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Discharge   bool      `json:"discharge"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"-"`
	BloodGroup  string    `json:"bloodGroup"`
	Description string    `json:"description"`
}
