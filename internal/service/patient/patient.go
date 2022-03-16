package service

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/anish-kmr/patient-system/internal/model"
	store "github.com/anish-kmr/patient-system/internal/store/patient"
)

type PatientService struct {
	store store.Patient
}

func New(s store.Patient) *PatientService {
	return &PatientService{store: s}
}

func (ps *PatientService) GetByID(ctx *gofr.Context, id int) (*model.Patient, error) {
	if id <= 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	patient, err := ps.store.GetByID(ctx, id)

	return patient, err
}

func (ps *PatientService) GetAll(ctx *gofr.Context, filters map[string]string) ([]*model.Patient, error) {
	patients, err := ps.store.GetAll(ctx, filters)
	return patients, err
}

func (ps *PatientService) Create(ctx *gofr.Context, patient *model.Patient) (*model.Patient, error) {
	if !validateName(patient.Name) {
		return nil, errors.InvalidParam{Param: []string{"name"}}
	}

	if !validatePhone(patient.Phone) {
		return nil, errors.InvalidParam{Param: []string{"phone"}}
	}

	if !validateBloodGroup(patient.BloodGroup) {
		return nil, errors.InvalidParam{Param: []string{"bloodgroup"}}
	}

	patient, err := ps.store.Create(ctx, patient)

	return patient, err
}
func (ps *PatientService) Update(ctx *gofr.Context, id int, patient *model.Patient) (*model.Patient, error) {
	if !patient.Name.IsZero() && !validateName(patient.Name) {
		return nil, errors.InvalidParam{Param: []string{"name"}}
	}

	if !patient.Phone.IsZero() && !validatePhone(patient.Phone) {
		return nil, errors.InvalidParam{Param: []string{"phone"}}
	}

	if !patient.BloodGroup.IsZero() && !validateBloodGroup(patient.BloodGroup) {
		return nil, errors.InvalidParam{Param: []string{"bloodgroup"}}
	}

	_, err := ps.store.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	updatedPatient, err := ps.store.Update(ctx, id, patient)

	return updatedPatient, err
}

func (ps *PatientService) Delete(ctx *gofr.Context, id int) error {
	if id <= 0 {
		return errors.InvalidParam{Param: []string{"id"}}
	}

	_, err := ps.store.GetByID(ctx, id)

	if err != nil {
		return err
	}

	return ps.store.Delete(ctx, id)
}
