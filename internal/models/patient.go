package models

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Patient struct {
	ID          int         `json:"id"`
	Name        null.String `json:"name"`
	Phone       null.String `json:"phone"`
	Discharged  null.Bool   `json:"discharged"`
	BloodGroup  null.String `json:"bloodGroup"`
	Description null.String `json:"description"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}
