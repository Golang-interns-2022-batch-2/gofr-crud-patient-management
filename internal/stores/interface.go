package stores

//import "github.com/aakanksha/patient-management/patient-management-system/internal/models"

//go:generate mockgen -destination=interface_mock.go -package=stores github.com/aakanksha/patient-management-system/internal/stores PatientInterface
//import "github.com/aakanksha/patient-management/patient-management-system/internal/models"
import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/aakanksha/updated-patient-management-system/internal/models"
)

type StoreInterface interface {
	GetByID(*gofr.Context, int) (*models.Patient, error)
	Insert(*gofr.Context, *models.Patient) (*models.Patient, error)
	Update(*gofr.Context, *models.Patient) (*models.Patient, error)
	Delete(*gofr.Context, int) error
	GetAll(*gofr.Context) ([]*models.Patient, error)
}
