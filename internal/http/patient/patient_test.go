package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gofrerrors "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	httperrors "github.com/anish-kmr/patient-system/internal/errors"
	"github.com/anish-kmr/patient-system/internal/model"
	service "github.com/anish-kmr/patient-system/internal/service/patient"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

var patient = model.Patient{
	ID:          1,
	Name:        null.StringFrom("Anish"),
	Phone:       null.StringFrom("+91 9999999999"),
	Discharged:  null.BoolFrom(false),
	BloodGroup:  null.StringFrom("+B"),
	Description: null.StringFrom("Stress"),
	CreatedAt:   null.TimeFrom(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)),
	UpdatedAt:   null.TimeFrom(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)),
}

func TestGetByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := service.NewMockPatient(mockCtrl)

	testcases := []struct {
		description string
		id          string
		resBody     interface{}
		mock        *gomock.Call
		err         error
	}{
		{
			description: "Success Case",
			id:          "1",
			resBody: httpResponse{
				Code:   http.StatusOK,
				Status: "SUCCESS",
				Data:   data{&patient},
			},
			err: nil,
			mock: mockService.EXPECT().GetByID(gomock.Any(), 1).Return(
				&model.Patient{
					ID:          1,
					Name:        null.StringFrom("Anish"),
					Phone:       null.StringFrom("+91 9999999999"),
					Discharged:  null.BoolFrom(false),
					BloodGroup:  null.StringFrom("+B"),
					Description: null.StringFrom("Stress"),
					CreatedAt:   null.TimeFrom(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)),
					UpdatedAt:   null.TimeFrom(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)),
				}, nil,
			),
		},
		{
			description: "invalid Id Param",
			id:          "asd",
			resBody:     nil,
			err:         gofrerrors.InvalidParam{Param: []string{"id"}},
		},
		{
			description: "Error From Service",
			id:          "1",
			resBody:     nil,
			mock: mockService.EXPECT().GetByID(gomock.Any(), 1).Return(
				nil, errors.New(httperrors.PatientNotFound),
			),
			err: errors.New(httperrors.PatientNotFound),
		},
	}
	handler := New(mockService)

	for i, tc := range testcases {
		r := httptest.NewRequest("GET", "/patient", nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		resp, err := handler.GetByID(ctx)
		assert.Equal(t, resp, tc.resBody, "[Test %v] : Expected :%v Got %v ", i, resp, tc.resBody)
		assert.Equal(t, err, tc.err, "[Test %v] : Expected :%v Got %v ", i, err, tc.err)
	}
}

func TestGetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := service.NewMockPatient(mockCtrl)
	queryFilters := map[string]string{"page": "0", "limit": "1"}
	testcases := []struct {
		description string
		filter      map[string]string
		resBody     interface{}
		mock        *gomock.Call
		err         error
	}{
		{
			description: "Success Case",
			filter:      map[string]string{"page": "1", "limit": "1"},
			resBody: httpResponse{
				Code:   http.StatusOK,
				Status: "SUCCESS",
				Data:   data{[]*model.Patient{&patient}},
			},
			mock: mockService.EXPECT().GetAll(gomock.Any(), queryFilters).Return(
				[]*model.Patient{
					{
						ID:          1,
						Name:        null.StringFrom("Anish"),
						Phone:       null.StringFrom("+91 9999999999"),
						Discharged:  null.BoolFrom(false),
						BloodGroup:  null.StringFrom("+B"),
						Description: null.StringFrom("Stress"),
						CreatedAt:   null.TimeFrom(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)),
						UpdatedAt:   null.TimeFrom(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)),
					},
				}, nil,
			),
		},
		{
			description: "Wrong Case",
			filter:      queryFilters,
			resBody:     nil,
			mock: mockService.EXPECT().GetAll(gomock.Any(), queryFilters).Return(
				nil, errors.New(httperrors.PatientFailed),
			),
			err: errors.New(httperrors.PatientFailed),
		},
	}
	handler := New(mockService)

	for i, tc := range testcases {
		r := httptest.NewRequest("GET", "/patient?page=0&limit=1", nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		resp, err := handler.GetAll(ctx)
		assert.Equal(t, resp, tc.resBody, "[Test %v] : Expected :%v Got %v ", i, resp, tc.resBody)
		assert.Equal(t, err, tc.err, "[Test %v] : Expected :%v Got %v ", i, err, tc.err)
	}
}

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := service.NewMockPatient(mockCtrl)

	testcases := []struct {
		description string
		reqBody     []byte
		resBody     interface{}
		err         error
		mock        *gomock.Call
	}{
		{
			description: "Response 200",
			reqBody:     []byte(`{"name": "Anish Kumar","phone": "+91 9999999999","bloodGroup": "+B","description": "Fever"}`),
			resBody: httpResponse{
				Code:   http.StatusOK,
				Status: "SUCCESS",
				Data:   data{&patient},
			},

			mock: mockService.
				EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(&patient, nil),
		},
		{
			description: "Error From Service",
			reqBody:     []byte(`{"name": "Anish Kumar","phone": "+91 9999999999","bloodGroup": "+B","description": "Fever"}`),
			resBody:     nil,
			mock: mockService.
				EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(nil, errors.New(httperrors.CreateFailed)),
			err: errors.New(httperrors.CreateFailed),
		},
		{
			description: "Response 200",
			reqBody:     []byte(`{"name": "Anish Kumar" "phone": "+91 9999999999","bloodGroup": "+B","description": "Fever"}`),
			resBody:     nil,
			err:         gofrerrors.InvalidParam{Param: []string{"id"}},
		},
	}
	handler := New(mockService)

	for i, tc := range testcases {
		r := httptest.NewRequest("POST", "/patient", bytes.NewReader(tc.reqBody))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		resp, err := handler.Create(ctx)
		assert.Equal(t, resp, tc.resBody, "[Test %v] : Expected :%v Got %v ", i, resp, tc.resBody)
		assert.Equal(t, err, tc.err, "[Test %v] : Expected :%v Got %v ", i, err, tc.err)
	}
}
func TestUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := service.NewMockPatient(mockCtrl)

	testcases := []struct {
		description string
		id          string
		reqBody     []byte
		resBody     interface{}
		mock        *gomock.Call
		err         error
	}{
		{
			description: "Response 200",
			id:          "1",
			reqBody: []byte(`{` +
				`"id":1,` +
				`"name":"Anish Kumar",` +
				`"phone":"+91 9999999999",` +
				`"discharged":null,` +
				`"bloodGroup":"+B",` +
				`"description":"Fever",` +
				`"createdAt":"2022-01-01T01:01:01.000000001Z",` +
				`"updatedAt":"2022-01-01T01:01:01.000000001Z"` +
				`}`),

			resBody: httpResponse{
				Code:   200,
				Status: "SUCCESS",
				Data:   data{&patient},
			},
			err: nil,
			mock: mockService.
				EXPECT().
				Update(gomock.Any(), 1, gomock.Any()).
				Return(&patient, nil),
		},
		{
			description: "Error From Service",
			id:          "1",
			reqBody: []byte(`{"name": "Anish Kumar",` +
				`"phone": "+91 9999999999",` +
				`"bloodGroup": "+B",` +
				`"description": "Fever"}`),
			resBody: nil,
			err:     errors.New(httperrors.UpdateFailed),
			mock: mockService.
				EXPECT().
				Update(gomock.Any(), 1, gomock.Any()).
				Return(nil, errors.New(httperrors.UpdateFailed)),
		},
		{
			description: "Parse Error",
			id:          "1",
			reqBody:     []byte(`{"name": "Anish Kumar" "phone": "+91 9999999999","bloodGroup": "+B","description": "Fever"}`),
			resBody:     nil,
			err:         gofrerrors.InvalidParam{Param: []string{"body"}},
		},
		{
			description: "Response 200",
			id:          "abc",
			resBody:     nil,
			err:         gofrerrors.InvalidParam{Param: []string{"id"}},
		},
	}
	handler := New(mockService)

	for i, tc := range testcases {
		r := httptest.NewRequest("PUT", "/patient", bytes.NewReader(tc.reqBody))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		resp, err := handler.Update(ctx)

		assert.Equal(t, resp, tc.resBody, "[Test %v] : Expected :%v Got %v ", i, resp, tc.resBody)
		assert.Equal(t, err, tc.err, "[Test %v] : Expected :%v Got %v ", i, err, tc.err)
	}
}

func TestDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := service.NewMockPatient(mockCtrl)

	testcases := []struct {
		description string
		id          string
		resBody     interface{}
		err         error
		mock        *gomock.Call
	}{
		{
			description: "Response 200",
			id:          "1",
			resBody: httpResponse{
				Code:   200,
				Status: "SUCCESS",
				Data:   "Patient Deleted Successfully",
			},
			mock: mockService.
				EXPECT().
				Delete(gomock.Any(), 1).
				Return(nil),
		},
		{
			description: "Invalid ID Param",
			id:          "abc",
			resBody:     nil,
			err:         gofrerrors.InvalidParam{Param: []string{"id"}},
		},
		{
			description: "Error From Service",
			id:          "1",
			resBody:     nil,
			err:         errors.New(httperrors.DeleteFailed),
			mock: mockService.
				EXPECT().
				Delete(gomock.Any(), 1).
				Return(errors.New(httperrors.DeleteFailed)),
		},
	}
	handler := New(mockService)

	for i, tc := range testcases {
		r := httptest.NewRequest("DELETE", "/patient", nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)
		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		resp, err := handler.Delete(ctx)

		assert.Equal(t, resp, tc.resBody, "[Test %v] : Expected :%v Got %v ", i, resp, tc.resBody)
		assert.Equal(t, err, tc.err, "[Test %v] : Expected :%v Got %v ", i, err, tc.err)
	}
}
