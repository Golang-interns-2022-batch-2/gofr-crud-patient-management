package store

import (
	"GOFR/models"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Patient interface {
	GetByID(ctx *gofr.Context, id int) (*models.Patient, error)
	Create(ctx *gofr.Context, patient *models.Patient) (*models.Patient, error)
	Get(ctx *gofr.Context) ([]*models.Patient, error)
	Update(ctx *gofr.Context, id int, patient *models.Patient) (*models.Patient, error)
	Delete(ctx *gofr.Context, id int) error
}
