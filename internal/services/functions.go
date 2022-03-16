package services

import (
	"regexp"
	"strings"

	"github.com/punitj12/patient-app-gofr/internal/models"
)

func validate(patient *models.Patient) bool {
	// Validate Name
	if patient.Name.Valid && patient.Name.String == "" {
		return false
	}

	// Validate Blood Group
	if patient.BloodGroup.Valid {
		bg := strings.ToUpper(patient.BloodGroup.String)
		switch bg {
		case "+AB", "-AB", "+A", "-A", "+B", "-B", "+O", "-O":

		default:
			return false
		}
	}

	// Validate Phone Number
	if patient.Phone.Valid {
		reg := regexp.MustCompile(`^(\+\d{1,2}\s?)?1?\-?\.?\s?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)
		if !reg.MatchString(patient.Phone.String) {
			return false
		}
	}

	return true
}
