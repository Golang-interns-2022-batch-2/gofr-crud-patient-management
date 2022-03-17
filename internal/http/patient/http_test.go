package patient

import (
	"bytes"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"errors"
	"github.com/aakanksha/updated-patient-management-system/internal/models"
	"github.com/aakanksha/updated-patient-management-system/internal/service"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

var patient = models.Patient{
	Id:          1,
	Name:        "ZopSmart",
	Phone:       "+919172681679",
	Discharge:   true,
	BloodGroup:  "+A",
	Description: "patient description",
}

func TestGetByID(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPatientService := service.NewMockServiceInterface(mockCtrl)
	app := gofr.New()
	testCases := []struct {
		id            string
		mockCall      *gomock.Call
		expectedError error
	}{

		{
			id:            "1",
			mockCall:      mockPatientService.EXPECT().GetByID(gomock.Any(), 1).Return(&models.Patient{}, nil),
			expectedError: nil,
		},

		{
			id:            "2",
			mockCall:      mockPatientService.EXPECT().GetByID(gomock.Any(), 2).Return(&patient, errors.New("error found")),
			expectedError: errors.New("error found"),
		},
	}
	p := New(mockPatientService)
	for _, testCase := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": testCase.id,
		})
		_, err := p.GetByID(ctx)
		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("expected error :%v, got :%v ", testCase.expectedError, err)
		}
	}
}

func Test_DeletePatientService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPatientService := service.NewMockServiceInterface(mockCtrl)
	app := gofr.New()
	testCases := []struct {
		body          []byte
		id            string
		mockCall      *gomock.Call
		expectedError error
		status        int
	}{
		// Success
		{
			id:            "1",
			mockCall:      mockPatientService.EXPECT().Delete(gomock.Any(), 1).Return(nil),
			expectedError: nil,
			//status:        200,
		},
		//Failure
		{
			id:            "2",
			mockCall:      mockPatientService.EXPECT().Delete(gomock.Any(), 2).Return(errors.New("error")),
			expectedError: errors.New("error"),
			//status:        500,
		},
	}
	p := New(mockPatientService)
	for _, testCase := range testCases {

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "http://localhost:8080", nil)
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": testCase.id,
		})
		_, err := p.Delete(ctx)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("expected error :%v, got :%v ", testCase.expectedError, err)
		}
	}
}

func Test_GetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPatientService := service.NewMockServiceInterface(mockCtrl)
	app := gofr.New()
	testCases := []struct {
		body          []byte
		mockCall      *gomock.Call
		expectedError error
		status        int
	}{
		// Success
		{
			mockCall:      mockPatientService.EXPECT().GetAll(gomock.Any()).Return([]*models.Patient{&patient}, nil),
			expectedError: nil,
			//	status:        200,
		},
		//Failure
		{
			mockCall:      mockPatientService.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("error")),
			expectedError: errors.New("error"),
			//	status:        500,
		},
	}
	p := New(mockPatientService)
	for _, testCase := range testCases {

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := p.GetAll(ctx)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("expected error :%v, got :%v ", testCase.expectedError, err)
		}

	}
}

func Test_Insert(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPatientService := service.NewMockServiceInterface(mockCtrl)
	app := gofr.New()
	testCases := []struct {
		body          []byte
		input         models.Patient
		mockCall      *gomock.Call
		expectedError error
		//status        int
	}{
		//Success
		{
			body: []byte(`{
				"name": "Zopsmart",
				"phone": "+919172681679",
				"discharge": true,
				"bloodGroup": "+A",
				"description": "patient description"
				}`),
			input:         patient,
			mockCall:      mockPatientService.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(&patient, nil),
			expectedError: nil,
			//status:        200,
		},
		{
			body: []byte(`{
                   "name":"Zopsmart",
                   "phone":"7676878978786",
                   "discharge":true,
                   "bloodgroup":"+A",
                   "description":"patient description"
                        }`),
			input:         patient,
			mockCall:      mockPatientService.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil, errors.New("invalid fields")),
			expectedError: errors.New("invalid fields"),
		},
	}
	p := New(mockPatientService)
	for _, testCase := range testCases {

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8080", bytes.NewReader(testCase.body))
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := p.Insert(ctx)

		if !reflect.DeepEqual(testCase.expectedError, err) {
			t.Errorf("expected error :%v, got :%v ", testCase.expectedError, err)
		}

	}
}

func Test_UpdatePatientService(t *testing.T) {

	var pat = &models.Patient{
		Name:        "ZopSmart",
		Description: "patient description",
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	idString := strconv.Itoa(patient.Id)
	mockPatientService := service.NewMockServiceInterface(mockCtrl)
	app := gofr.New()
	testCases := []struct {
		body          []byte
		id            string
		mockCall      *gomock.Call
		expectedError error
		//status        int
	}{
		// Success
		{
			body: []byte(`{
				"name": "ZopSmart",
				"description": "patient description"
				}`),
			id:            idString,
			mockCall:      mockPatientService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(pat, nil),
			expectedError: nil,
			//status:        200,
		},
		{
			body: []byte(`{
				"name": 1,
				"description": "patient description"
				}`),
			expectedError: errors.New("cannot read from body"),
		},
	}
	p := New(mockPatientService)
	for _, testCase := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "http://localhost:8080", bytes.NewReader(testCase.body))
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := p.Update(ctx)

		if !reflect.DeepEqual(testCase.expectedError, err) {
			t.Errorf("expected error :%v, got :%v ", testCase.expectedError, err)
		}
	}
}
