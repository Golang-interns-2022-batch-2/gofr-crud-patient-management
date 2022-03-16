package patient

import (
	"net/http"
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"github.com/shivanisharma200/patient-management/internal/models"
	"github.com/shivanisharma200/patient-management/internal/service"
)

type API struct {
	PatientService service.Patient
}

func New(patientService service.Patient) *API {
	return &API{PatientService: patientService}
}

type ErrorStorer struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"Message"`
}
type ResponseStorer struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type data struct {
	Patient interface{}
}

// GetPatient
func (p *API) GetByID(ctx *gofr.Context) (interface{}, error) {
	var response interface{}

	idString := ctx.PathParam("id")
	id, _ := strconv.Atoi(idString)

	patient, err := p.PatientService.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	response = ResponseStorer{
		Code:   http.StatusOK,
		Status: "SUCCESS",
		Data:   data{patient},
	}
	r := types.Response{
		Data: response,
	}

	return r, nil
}

// createPatient
func (p *API) Create(ctx *gofr.Context) (interface{}, error) {
	var response interface{}

	var patient models.Patient

	err := ctx.Bind(&patient)
	if err != nil {
		return nil, errors.Error("cannot read from body")
	}

	patientVal, err := p.PatientService.Create(ctx, &patient)

	if err != nil {
		return nil, err
	}

	patientVal = &models.Patient{ID: patientVal.ID, Name: patientVal.Name, Phone: patientVal.Phone,
		Discharged: patientVal.Discharged, CreatedAt: patientVal.CreatedAt, UpdatedAt: patientVal.UpdatedAt,
		BloodGroup: patientVal.BloodGroup, Description: patientVal.Description}
	response = ResponseStorer{
		Code:   http.StatusOK,
		Status: "SUCCESS",
		Data:   data{patientVal},
	}
	r := types.Response{
		Data: response,
	}

	return r, nil
}

// updatePatient
func (p *API) Update(ctx *gofr.Context) (interface{}, error) {
	var response interface{}

	idString := ctx.PathParam("id")
	id, _ := strconv.Atoi(idString)

	var patient *models.Patient

	err := ctx.Bind(&patient)
	if err != nil {
		return nil, errors.Error("cannot read from body")
	}

	patient, err = p.PatientService.Update(ctx, id, patient)

	if err != nil {
		return nil, err
	}

	response = ResponseStorer{
		Code:   http.StatusOK,
		Status: "SUCCESS",
		Data:   data{patient},
	}

	r := types.Response{
		Data: response,
	}

	return r, nil
}

// GetPatients
func (p *API) Get(ctx *gofr.Context) (interface{}, error) {
	var response interface{}

	patients, err := p.PatientService.Get(ctx)

	if err != nil {
		return nil, err
	}

	response = ResponseStorer{
		Code:   http.StatusOK,
		Status: "SUCCESS",
		Data:   data{patients},
	}

	r := types.Response{
		Data: response,
	}

	return r, nil
}

// deletePatient
func (p *API) Delete(ctx *gofr.Context) (interface{}, error) {
	var response interface{}

	idString := ctx.PathParam("id")
	id, _ := strconv.Atoi(idString)
	err := p.PatientService.Delete(ctx, id)

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
