package store

//go:generate mockgen -destination=interface_mock.go -package=store github.com/shivanisharma200/patient-management/internal/store Patient

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/shivanisharma200/patient-management/internal/models"
)

type Patient interface {
	GetByID(ctx *gofr.Context, id int) (*models.Patient, error)
	Create(ctx *gofr.Context, patient *models.Patient) (*models.Patient, error)
	Get(ctx *gofr.Context) ([]*models.Patient, error)
	Update(ctx *gofr.Context, id int, patient *models.Patient) (*models.Patient, error)
	Delete(ctx *gofr.Context, id int) error
}
