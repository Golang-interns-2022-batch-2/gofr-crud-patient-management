package http

import (
	"fmt"
	"net/http"
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/anish-kmr/patient-system/internal/model"
	service "github.com/anish-kmr/patient-system/internal/service/patient"
)

type PatientHandler struct {
	service service.Patient
}

type data struct {
	Patient interface{}
}

type httpResponse struct {
	Code   int
	Status string
	Data   interface{}
}

func New(ps service.Patient) *PatientHandler {
	return &PatientHandler{service: ps}
}

func (ph *PatientHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	var response interface{}

	param := ctx.PathParam("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	patient, err := ph.service.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response = httpResponse{
		Code:   http.StatusOK,
		Status: "SUCCESS",
		Data:   data{patient},
	}

	return response, nil
}

func (ph *PatientHandler) GetAll(ctx *gofr.Context) (interface{}, error) {
	var response interface{}

	filters := ctx.Params()

	patients, err := ph.service.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	response = httpResponse{
		Code:   http.StatusOK,
		Status: "SUCCESS",
		Data:   data{patients},
	}

	return response, nil
}

func (ph *PatientHandler) Create(ctx *gofr.Context) (interface{}, error) {
	var response interface{}

	var patient model.Patient
	err := ctx.Bind(&patient)
	fmt.Println("BODY ", patient)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	newPatient, err := ph.service.Create(ctx, &patient)
	if err != nil {
		return nil, err
	}

	response = httpResponse{
		Code:   http.StatusOK,
		Status: "SUCCESS",
		Data:   data{newPatient},
	}

	return response, nil
}

func (ph *PatientHandler) Update(ctx *gofr.Context) (interface{}, error) {
	var response interface{}

	params := ctx.PathParam("id")
	id, err := strconv.Atoi(params)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	var patient model.Patient
	err = ctx.Bind(&patient)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	updatedPatient, err := ph.service.Update(ctx, id, &patient)
	if err != nil {
		return nil, err
	}

	response = httpResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   data{updatedPatient},
	}

	return response, nil
}

func (ph *PatientHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	var response interface{}

	param := ctx.PathParam("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	err = ph.service.Delete(ctx, id)

	if err != nil {
		return nil, err
	}

	response = httpResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   "Patient Deleted Successfully",
	}

	return response, nil
}
