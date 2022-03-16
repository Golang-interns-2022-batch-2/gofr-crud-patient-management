package patient

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"GOFR/models"
	"GOFR/store"
)

type Patient struct {
	PatientStoreHandler store.Patient
}

func New(str store.Patient) *Patient {
	return &Patient{PatientStoreHandler: str}
}


// GetPatientService
func (p *Patient) GetByID(ctx *gofr.Context, idString string) (*models.Patient, error) {
	id, _ := strconv.Atoi(idString)
	if !IsIDValid(id) {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	return p.PatientStoreHandler.GetByID(ctx, id)
}

// CreatePatientService
func (p *Patient) Create(ctx *gofr.Context, patient *models.Patient) (*models.Patient, error) {
	err := patient.Validate()

	if err != nil {
		return nil, errors.Error("invalid fileds")
	}

	return p.PatientStoreHandler.Create(ctx, patient)
}

// GetAllService
func (p *Patient) Get(ctx *gofr.Context) ([]*models.Patient, error) {
	return p.PatientStoreHandler.Get(ctx)
}

// UpdatePatientService
func (p *Patient) Update(ctx *gofr.Context, idString string, patient *models.Patient) (*models.Patient, error) {
	id, _ := strconv.Atoi(idString)
	if !IsIDValid(id) {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	_, err := p.GetByID(ctx, idString)

	if err != nil {
		return nil, errors.EntityNotFound{Entity: "Patient", ID: "id"}
	}

	return p.PatientStoreHandler.Update(ctx, id, patient)
}

// DeletePatientService
func (p *Patient) Delete(ctx *gofr.Context, idString string) error {
	id, _ := strconv.Atoi(idString)
	if !IsIDValid(id) {
		return errors.InvalidParam{Param: []string{"id"}}
	}

	_, err := p.GetByID(ctx, idString)

	if err != nil {
		return errors.InvalidParam{Param: []string{"id"}}
	}

	return p.PatientStoreHandler.Delete(ctx, id)
}
