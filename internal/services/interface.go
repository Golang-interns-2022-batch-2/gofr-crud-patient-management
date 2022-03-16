//go:generate mockgen -destination=interface_mock.go -package=services github.com/punitj12/patient-app-gofr/internal/services PatientServicer
//nolint
package services

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/punitj12/patient-app-gofr/internal/models"
)

type PatientServicer interface {
	Create(*gofr.Context, *models.Patient) (*models.Patient, error)
	Delete(*gofr.Context, int) error
	Get(*gofr.Context, int) (*models.Patient, error)
	GetAll(*gofr.Context) ([]*models.Patient, error)
	Update(*gofr.Context, *models.Patient) (*models.Patient, error)
}
