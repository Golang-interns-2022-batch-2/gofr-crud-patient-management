package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"github.com/golang/mock/gomock"
	"github.com/punitj12/patient-app-gofr/internal/models"
	"github.com/punitj12/patient-app-gofr/internal/services"
	"gopkg.in/guregu/null.v4"
)

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pat := models.Patient{
		ID:          1,
		Name:        null.StringFrom("Abhishek"),
		Phone:       null.StringFrom("+919232323232"),
		BloodGroup:  null.StringFrom("+AB"),
		Description: null.StringFrom("Head Ache"),
	}
	errPat := models.Patient{
		ID:          1,
		Name:        null.StringFrom(""),
		Phone:       null.StringFrom("+919232323232"),
		BloodGroup:  null.StringFrom("+AB"),
		Description: null.StringFrom("Head Ache"),
	}
	PatientServiceMock := services.NewMockPatientServicer(mockCtrl)

	tcs := []struct {
		desc       string
		id         int
		body       []byte
		output     interface{}
		outputCode int
		mock       *gomock.Call
		err        error
	}{

		{
			desc:       "success case",
			id:         1,
			body:       []byte(`{"id":1,"name":"Abhishek","bloodGroup":"+AB","description":"Head Ache","phone":"+919232323232"}`),
			outputCode: 200,
			output: types.Response{
				Data: res{
					Data: &models.Patient{
						ID:          1,
						Name:        null.StringFrom("Abhishek"),
						Phone:       null.StringFrom("+919232323232"),
						Discharged:  null.BoolFrom(false),
						BloodGroup:  null.StringFrom("+AB"),
						Description: null.StringFrom("Head Ache"),
					},
				},
			},
			mock: PatientServiceMock.
				EXPECT().
				Create(gomock.Any(), &pat).
				Return(&models.Patient{
					ID:          1,
					Name:        null.StringFrom("Abhishek"),
					Phone:       null.StringFrom("+919232323232"),
					Discharged:  null.BoolFrom(false),
					BloodGroup:  null.StringFrom("+AB"),
					Description: null.StringFrom("Head Ache"),
				}, nil),
		},
		{
			desc:       "incorrect data format",
			id:         2,
			body:       []byte(`{"id":1,"name":"","bloodGroup":"+AB","description":"Head Ache","phone":"+919232323232"}`),
			outputCode: 400,
			mock: PatientServiceMock.
				EXPECT().
				Create(gomock.Any(), &errPat).
				Return(nil, errors.InvalidParam{}),
			err: errors.InvalidParam{},
		},
		{
			desc:       "error parsing body",
			id:         3,
			body:       []byte(`{"id":1"name":"","bloodGroup":"+AB","description":"Head Ache","phone":"+919232323232"}`),
			outputCode: 400,
			err:        errors.InvalidParam{},
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run("testing adding service", func(t *testing.T) {
			r := httptest.NewRequest("POST", "/patients", bytes.NewReader(tc.body))
			wr := httptest.NewRecorder()

			req := request.NewHTTPRequest(r)
			w := responder.NewContextualResponder(wr, r)
			ctx := gofr.NewContext(w, req, gofr.New())
			ctx.Context = context.Background()

			a := New(PatientServiceMock)
			res, er := a.Create(ctx)
			if er != nil && (tc.err.Error() != er.Error()) {
				t.Errorf("desc ->>> %v Expected Error : %v Got : %v", tc.desc, tc.err, er)
			}

			if !reflect.DeepEqual(res, tc.output) {
				t.Errorf("desc ->>> %v Expected Result : %v Got : %v", tc.desc, tc.output, res)
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	parsedTime, _ := time.Parse(time.RFC1123Z, "2022-03-05T18:44:05Z")
	pat := models.Patient{
		ID:          1,
		Name:        null.StringFrom("Abhi"),
		Phone:       null.StringFrom("+916666666666"),
		Discharged:  null.BoolFrom(false),
		BloodGroup:  null.StringFrom("+AB"),
		Description: null.StringFrom("Headache"),
		CreatedAt:   parsedTime,
		UpdatedAt:   parsedTime,
	}
	PatientServiceMock := services.NewMockPatientServicer(mockCtrl)

	tcs := []struct {
		desc       string
		id         string
		body       []byte
		output     interface{}
		outputCode int
		mock       *gomock.Call
		err        error
	}{
		{
			desc: "success case",
			id:   "1",
			output: types.Response{
				Data: res{
					Data: &pat,
				},
			},
			outputCode: 200,
			mock: PatientServiceMock.EXPECT().Get(gomock.Any(), 1).
				Return(&pat, nil),
		},
		{
			desc:       "invalid patient",
			id:         "2",
			outputCode: 404,
			mock: PatientServiceMock.EXPECT().Get(gomock.Any(), 2).
				Return(nil, errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(2)}),
			err: errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(2)},
		},
		{
			desc:       "other error",
			id:         "3",
			outputCode: 500,
			mock: PatientServiceMock.EXPECT().Get(gomock.Any(), 3).
				Return(&models.Patient{}, errors.Error("error fetching patient")),
			err: errors.Error("error fetching patient"),
		},
		{
			desc:       "invalid id format",
			id:         "a",
			outputCode: 400,
			err:        errors.InvalidParam{Param: []string{"id"}},
		},
	}

	for _, tc := range tcs {
		tc := tc

		r := httptest.NewRequest("GET", "/patients", bytes.NewReader(tc.body))
		wr := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		w := responder.NewContextualResponder(wr, r)

		ctx := gofr.NewContext(w, req, gofr.New())
		ctx.SetPathParams(map[string]string{"id": tc.id})

		a := New(PatientServiceMock)
		res, er := a.Get(ctx)

		if !reflect.DeepEqual(res, tc.output) {
			t.Errorf("desc ->>> %v Expected Result : %v Got : %v", tc.desc, tc.output, res)
		}

		if er != nil && (tc.err.Error() != er.Error()) {
			t.Errorf("desc ->>> %v Expected Error : %v Got : %v", tc.desc, tc.err, er)
		}
	}
}

func TestGetAllHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	parsedTime, _ := time.Parse(time.RFC1123Z, "2022-03-05T18:44:05Z")
	pat := []*models.Patient{
		{
			ID:          1,
			Name:        null.StringFrom("Abhi"),
			Phone:       null.StringFrom("+916666666666"),
			Discharged:  null.BoolFrom(false),
			BloodGroup:  null.StringFrom("+AB"),
			Description: null.StringFrom("Headache"),
			CreatedAt:   parsedTime,
			UpdatedAt:   parsedTime,
		},
	}
	PatientServiceMock := services.NewMockPatientServicer(mockCtrl)

	tcs := []struct {
		desc       string
		id         int
		body       []byte
		output     interface{}
		outputCode int
		mock       *gomock.Call
		err        error
	}{
		{
			desc: "success case",
			id:   1,
			output: types.Response{
				Data: res{
					Data: pat,
				},
			},
			outputCode: 200,
			mock: PatientServiceMock.EXPECT().GetAll(gomock.Any()).
				Return(pat, nil),
		},
		{
			desc:       "db error",
			id:         1,
			outputCode: 500,
			mock: PatientServiceMock.EXPECT().GetAll(gomock.Any()).
				Return(nil, errors.Error("error fetching patient data")),
			err: errors.Error("error fetching patient data"),
		},
		{
			desc:       "no data",
			id:         1,
			outputCode: 500,
			mock: PatientServiceMock.EXPECT().GetAll(gomock.Any()).
				Return(nil, errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(1)}),
			err: errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(1)},
		},
	}
	for _, tc := range tcs {
		tc := tc

		t.Run("testing get add handler", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/patients", bytes.NewReader(tc.body))
			wr := httptest.NewRecorder()

			req := request.NewHTTPRequest(r)
			w := responder.NewContextualResponder(wr, r)

			ctx := gofr.NewContext(w, req, gofr.New())

			a := New(PatientServiceMock)
			res, er := a.GetAll(ctx)

			if er != nil && (tc.err.Error() != er.Error()) {
				t.Errorf("desc ->>> %v Expected Error : %v Got : %v", tc.desc, tc.err, er)
			}

			if !reflect.DeepEqual(res, tc.output) {
				t.Errorf("desc ->>> %v Expected Result : %v Got : %v", tc.desc, tc.output, res)
			}
		})
	}
}

func TestUpdateHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pat := models.Patient{
		ID:          1,
		Name:        null.StringFrom("Abhishek"),
		Phone:       null.StringFrom("+919232323232"),
		BloodGroup:  null.StringFrom("+AB"),
		Description: null.StringFrom("Head Ache"),
	}
	errPat := models.Patient{
		ID:          1,
		Name:        null.StringFrom(""),
		Phone:       null.StringFrom("+919232323232"),
		BloodGroup:  null.StringFrom("+AB"),
		Description: null.StringFrom("Head Ache"),
	}
	PatientServiceMock := services.NewMockPatientServicer(mockCtrl)

	tcs := []struct {
		desc       string
		id         string
		body       []byte
		output     interface{}
		outputCode int
		mock       *gomock.Call
		err        error
	}{

		{
			desc:       "success case",
			id:         "1",
			body:       []byte(`{"id":1,"name":"Abhishek","bloodGroup":"+AB","description":"Head Ache","phone":"+919232323232"}`),
			outputCode: 200,
			output: types.Response{
				Data: res{
					Data: &pat,
				},
			},
			mock: PatientServiceMock.
				EXPECT().
				Update(gomock.Any(), &pat).
				Return(&pat, nil),
		},
		{
			desc:       "incorrect data format",
			id:         "1",
			body:       []byte(`{"id":1,"name":"","bloodGroup":"+AB","description":"Head Ache","phone":"+919232323232"}`),
			outputCode: 400,
			mock: PatientServiceMock.
				EXPECT().
				Update(gomock.Any(), &errPat).
				Return(nil, errors.InvalidParam{}),
			err: errors.InvalidParam{},
		},
		{
			desc:       "error parsing body",
			id:         "1",
			body:       []byte(`{"id":1"name":"","bloodGroup":"+AB","description":"Head Ache","phone":"+919232323232"}`),
			outputCode: 400,
			err:        errors.InvalidParam{},
		},
		{
			desc:       "database error",
			id:         "1",
			body:       []byte(`{"id":1,"name":"Abhishek","bloodGroup":"+AB","description":"Head Ache","phone":"+919232323232"}`),
			outputCode: 500,
			err:        errors.Error("error adding patient"),
			mock: PatientServiceMock.
				EXPECT().
				Update(gomock.Any(), &pat).
				Return(&models.Patient{}, errors.Error("error adding patient")),
		},
		{
			desc:       "error parsing id",
			id:         "a",
			body:       []byte(`{"id":a,"name":"Punit","bloodGroup":"+AB","description":"Head Ache","phone":"+919232323232"}`),
			outputCode: 400,
			err:        errors.InvalidParam{Param: []string{"id"}},
		},
	}

	for _, tc := range tcs {
		tc := tc

		r := httptest.NewRequest("PUT", "/patients", bytes.NewReader(tc.body))
		wr := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		w := responder.NewContextualResponder(wr, r)

		ctx := gofr.NewContext(w, req, gofr.New())
		ctx.SetPathParams(map[string]string{"id": tc.id})

		a := New(PatientServiceMock)
		res, er := a.Update(ctx)

		if er != nil && (tc.err.Error() != er.Error()) {
			t.Errorf("desc ->>> %v Expected Error : %v Got : %v", tc.desc, tc.err, er)
		}

		if !reflect.DeepEqual(res, tc.output) {
			t.Errorf("desc ->>> %v Expected Result : %v Got : %v", tc.desc, tc.output, res)
		}
	}
}

func TestDeleteHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	PatientServiceMock := services.NewMockPatientServicer(mockCtrl)

	tcs := []struct {
		desc       string
		id         string
		body       []byte
		output     interface{}
		outputCode int
		mock       *gomock.Call
		err        error
	}{
		{
			desc: "success case",
			id:   "1",
			output: types.Response{
				Data: res{
					Data: "Patient deleted successfully",
				},
			},
			outputCode: 200,
			mock: PatientServiceMock.EXPECT().Delete(gomock.Any(), 1).
				Return(nil),
			err: nil,
		},
		{
			desc:       "patient not found",
			id:         "2",
			outputCode: 404,
			mock: PatientServiceMock.EXPECT().Delete(gomock.Any(), 2).
				Return(errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(2)}),
			err: errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(2)},
		},
		{
			desc:       "err deleting",
			id:         "3",
			outputCode: 500,
			mock: PatientServiceMock.EXPECT().Delete(gomock.Any(), 3).
				Return(errors.Error("error deleting patient")),
			err: errors.Error("error deleting patient"),
		},
		{
			desc:       "invalid id format",
			id:         "a",
			outputCode: 400,
			err:        errors.InvalidParam{Param: []string{"id"}},
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run("testing get add handler", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/patients", bytes.NewReader(tc.body))
			wr := httptest.NewRecorder()

			req := request.NewHTTPRequest(r)
			w := responder.NewContextualResponder(wr, r)

			ctx := gofr.NewContext(w, req, gofr.New())
			ctx.SetPathParams(map[string]string{"id": tc.id})

			a := New(PatientServiceMock)
			res, er := a.Delete(ctx)

			if er != nil && (tc.err.Error() != er.Error()) {
				t.Errorf("desc ->>> %v Expected Error : %v Got : %v", tc.desc, tc.err, er)
			}
			if !reflect.DeepEqual(res, tc.output) {
				t.Errorf("desc ->>> %v Expected Result : %v Got : %v", tc.desc, tc.output, res)
			}
		})
	}
}
