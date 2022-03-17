package service

//go:generate mockgen -destination=serviceinterface_mock.go -package=service github.com/aakanksha/patient-management-system/internal/service PatientInterface
import (
	//"github.com/aakanksha/patient-management/patient-management-system/internal/models"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/aakanksha/updated-patient-management-system/internal/models"
)

type ServiceInterface interface {
	GetByID(*gofr.Context, int) (*models.Patient, error)
	Insert(*gofr.Context, *models.Patient) (*models.Patient, error)
	Update(*gofr.Context, *models.Patient) (*models.Patient, error)
	Delete(*gofr.Context, int) error
	GetAll(s *gofr.Context) ([]*models.Patient, error)
}
