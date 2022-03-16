package service

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/shivanisharma200/patient-management/internal/models"
)

type Patient interface {
	GetByID(ctx *gofr.Context, id string) (*models.Patient, error)
	Create(ctx *gofr.Context, patient *models.Patient) (*models.Patient, error)
	Get(ctx *gofr.Context) ([]*models.Patient, error)
	Update(ctx *gofr.Context, id string, patient *models.Patient) (*models.Patient, error)
	Delete(ctx *gofr.Context, id string) error
}
