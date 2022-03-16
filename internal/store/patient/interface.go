//go:generate mockgen -destination=mock_interface.go -package=store github.com/anish-kmr/patient-system/internal/store/patient Patient
package store

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/anish-kmr/patient-system/internal/model"
)

type Patient interface {
	GetByID(*gofr.Context, int) (*model.Patient, error)

	GetAll(*gofr.Context, map[string]string) ([]*model.Patient, error)

	Create(*gofr.Context, *model.Patient) (*model.Patient, error)

	Update(*gofr.Context, int, *model.Patient) (*model.Patient, error)

	Delete(*gofr.Context, int) error
}
