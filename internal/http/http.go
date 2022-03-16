package http

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"github.com/punitj12/patient-app-gofr/internal/models"
	"github.com/punitj12/patient-app-gofr/internal/services"
)

type res struct {
	Data interface{} `json:"patient"`
}

type Handler struct {
	service services.PatientServicer
}

func New(s services.PatientServicer) *Handler {
	return &Handler{s}
}

func (h *Handler) Create(ctx *gofr.Context) (interface{}, error) {
	var patient models.Patient

	e := ctx.Bind(&patient)

	if e != nil {
		return nil, errors.InvalidParam{}
	}

	pat, err := h.service.Create(ctx, &patient)

	if err != nil {
		return nil, err
	}

	r := res{
		Data: pat,
	}

	return types.Response{
		Data: r,
	}, nil
}

func (h *Handler) Get(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	id, er := strconv.Atoi(i)

	if er != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	patient, err := h.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	r := res{
		Data: patient,
	}

	return types.Response{
		Data: r,
	}, nil
}

func (h *Handler) GetAll(ctx *gofr.Context) (interface{}, error) {
	patient, err := h.service.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	r := res{
		Data: patient,
	}

	return types.Response{
		Data: r,
	}, nil
}

func (h *Handler) Update(ctx *gofr.Context) (interface{}, error) {
	var patient models.Patient

	i := ctx.PathParam("id")
	id, er := strconv.Atoi(i)

	if er != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	e := ctx.Bind(&patient)
	if e != nil {
		return nil, errors.InvalidParam{}
	}

	patient.ID = id

	pat, err := h.service.Update(ctx, &patient)
	if err != nil {
		return nil, err
	}

	r := res{
		Data: pat,
	}

	return types.Response{
		Data: r,
	}, nil
}

func (h *Handler) Delete(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	id, er := strconv.Atoi(i)

	if er != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	err := h.service.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	r := res{
		Data: "Patient deleted successfully",
	}

	return types.Response{
		Data: r,
	}, nil
}
