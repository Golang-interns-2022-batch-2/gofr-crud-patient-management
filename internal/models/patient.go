package models

import "developer.zopsmart.com/go/gofr/pkg/errors"

type Patient struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Discharged  bool   `json:"discharged"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	DeletedAt   string `json:"-"`
	BloodGroup  string `json:"bloodGroup"`
	Description string `json:"description"`
}

func (p *Patient) Validate() error {
	if p.ID < 0 {
		return errors.Error("invalid fileds")
	}

	if p.Name == "" {
		return errors.Error("invalid fileds")
	}

	if len(p.Phone) != 13 || p.Phone[:3] != "+91" {
		return errors.Error("invalid fileds")
	}

	return nil
}
