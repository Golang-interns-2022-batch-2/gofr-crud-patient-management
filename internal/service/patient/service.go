package patient

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"errors"
	"fmt"
	"github.com/aakanksha/updated-patient-management-system/internal/models"
	"github.com/aakanksha/updated-patient-management-system/internal/stores"
)

type Svc struct {
	psvc stores.StoreInterface
}

func New(st stores.StoreInterface) *Svc {
	return &Svc{st}
}

func (p Svc) GetByID(c *gofr.Context, id int) (*models.Patient, error) {
	patient, err := p.psvc.GetByID(c, id)

	if err != nil {
		return nil, err
	}

	return patient, err
}

func (ps Svc) Insert(c *gofr.Context, p *models.Patient) (*models.Patient, error) {
	if !validatename(p.Name) {
		return &models.Patient{}, errors.New("invalid name")
	}
	patient, err := ps.psvc.Insert(c, p)
	fmt.Println("service error", err)
	if err != nil {
		return nil, err
	}

	patientgetbyid, err := ps.psvc.GetByID(c, patient.Id)

	return patientgetbyid, err
}

func (pp Svc) Update(c *gofr.Context, p *models.Patient) (*models.Patient, error) {
	_, err := pp.psvc.Update(c, p)
	if err != nil {
		fmt.Println("update serv")
		return nil, err
	}

	patient, err := pp.psvc.GetByID(c, p.Id)

	return patient, err
}

func (p Svc) Delete(c *gofr.Context, id int) error {
	err := p.psvc.Delete(c, id)
	if err != nil {
		return errors.New("unable to delete user")
	}

	return err
}

func (p Svc) GetAll(c *gofr.Context) ([]*models.Patient, error) {
	var patients []*models.Patient
	patients, err := p.psvc.GetAll(c)

	return patients, err
}
