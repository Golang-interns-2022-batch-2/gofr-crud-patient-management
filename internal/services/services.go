package services

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/punitj12/patient-app-gofr/internal/models"
	"github.com/punitj12/patient-app-gofr/internal/stores"
)

type Service struct {
	store stores.PatientStorer
}

func New(s stores.PatientStorer) *Service {
	return &Service{s}
}

func (s *Service) Create(ctx *gofr.Context, patient *models.Patient) (*models.Patient, error) {
	if !validate(patient) {
		return nil, errors.InvalidParam{}
	}

	patientNew, err := s.store.Create(ctx, patient)

	if err != nil {
		return nil, errors.Error("error adding patient")
	}

	return patientNew, nil
}

func (s *Service) Delete(ctx *gofr.Context, id int) error {
	if id <= 0 {
		return errors.InvalidParam{Param: []string{"id"}}
	}

	err := s.store.Delete(ctx, id)

	return err
}

func (s *Service) Get(ctx *gofr.Context, id int) (*models.Patient, error) {
	if id <= 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	patient, err := s.store.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	return patient, nil
}

func (s *Service) GetAll(ctx *gofr.Context) ([]*models.Patient, error) {
	patients, err := s.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (s *Service) Update(ctx *gofr.Context, patient *models.Patient) (*models.Patient, error) {
	if patient.ID <= 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	if !validate(patient) {
		return nil, errors.InvalidParam{}
	}

	patientNew, err := s.store.Update(ctx, patient)

	if err != nil {
		return nil, errors.Error("error updating patient")
	}

	return patientNew, nil
}
