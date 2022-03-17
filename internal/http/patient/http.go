package patient

import (
	//"github.com/aakanksha/updated-patient-management-system/internal/stores"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"errors"
	"fmt"
	"github.com/aakanksha/updated-patient-management-system/internal/models"
	"github.com/aakanksha/updated-patient-management-system/internal/service"
	"net/http"
	"strconv"
	//"strconv"
)

type PatientHandler struct {
	psvc service.ServiceInterface
}

func New(s service.ServiceInterface) *PatientHandler {
	return &PatientHandler{s}
}

type ResponseStorer struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type data struct {
	Patient interface{} `json:"patient"`
}

func (p *PatientHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	id, _ := strconv.Atoi(i)
	fmt.Println("HTTP LAYER")
	patient, err := p.psvc.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}
	Data := data{patient}
	r := types.Response{
		Data: Data,
	}
	return r, nil
}

func (p *PatientHandler) Insert(ctx *gofr.Context) (interface{}, error) {
	var pt models.Patient
	err := ctx.Bind(&pt)
	if err != nil {
		return nil, err
	}
	insertpatient, err := p.psvc.Insert(ctx, &pt)

	if err != nil {
		return nil, err
	}
	Data := data{insertpatient}
	r := types.Response{
		Data: Data,
	}
	return r, nil
}

func (p *PatientHandler) Update(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	id, _ := strconv.Atoi(i)
	var pt *models.Patient
	err := ctx.Bind(&pt)
	pt.Id = id
	if err != nil {
		return nil, errors.New("cannot read from body")
	}
	updatepatient, err := p.psvc.Update(ctx, pt)
	if err != nil {
		return &models.Patient{}, err
	}
	Data := data{updatepatient}
	r := types.Response{
		Data: Data,
	}
	return r, nil
}

func (p *PatientHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	var response interface{}
	i := ctx.PathParam("id")
	id, _ := strconv.Atoi(i)
	err := p.psvc.Delete(ctx, id)
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

func (p *PatientHandler) GetAll(ctx *gofr.Context) (interface{}, error) {

	patients, err := p.psvc.GetAll(ctx)
	if err != nil {
		return err, err
	}
	Data := data{patients}
	r := types.Response{
		Data: Data,
	}
	return r, err
}
