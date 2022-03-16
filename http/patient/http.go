package patient

import (
	"net/http"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"GOFR/models"
	"GOFR/service"
)

type Store struct {
	PatientService service.Patient
}

func New(patientService service.Patient) *Store {
	return &Store{PatientService: patientService}
}

type ResponseStorer struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type data struct {
	Patient interface{} `json:"patient"`
}

// GetPatient
func (p *Store) GetByID(ctx *gofr.Context) (interface{}, error) {
	idString := ctx.PathParam("id")

	patient, err := p.PatientService.GetByID(ctx, idString)

	if err != nil {
		return nil, err
	}

	Data := data{patient}
	r := types.Response{
		Data: Data,
	}

	return r, nil
}

// createPatient
func (p *Store) Create(ctx *gofr.Context) (interface{}, error) {
	var patient models.Patient

	err := ctx.Bind(&patient)
	if err != nil {
		return nil, errors.Error("cannot read from body")
	}

	patientVal, err := p.PatientService.Create(ctx, &patient)

	if err != nil {
		return nil, err
	}

	Data := data{patientVal}
	r := types.Response{
		Data: Data,
	}

	return r, nil
}

// updatePatient
func (p *Store) Update(ctx *gofr.Context) (interface{}, error) {
	idString := ctx.PathParam("id")

	var patient *models.Patient

	err := ctx.Bind(&patient)
	if err != nil {
		return nil, errors.Error("cannot read from body")
	}

	patient, err = p.PatientService.Update(ctx, idString, patient)

	if err != nil {
		return nil, err
	}

	Data := data{patient}
	r := types.Response{
		Data: Data,
	}

	return r, nil
}

// GetPatients
func (p *Store) Get(ctx *gofr.Context) (interface{}, error) {
	patients, err := p.PatientService.Get(ctx)

	if err != nil {
		return nil, err
	}

	Data := data{patients}
	r := types.Response{
		Data: Data,
	}

	return r, nil
}

// deletePatient
func (p *Store) Delete(ctx *gofr.Context) (interface{}, error) {
	var response interface{}

	idString := ctx.PathParam("id")
	err := p.PatientService.Delete(ctx, idString)

	if err != nil {
		return nil, err
	}

	response = ResponseStorer{
		Code:   http.StatusOK,
		Status: "SUCCESS",
		Data:   "Patient Deleted Successfully",
	}

	r := types.Response{
		Data: response,
	}

	return r, nil
}
