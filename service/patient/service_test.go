package patient

import (
	"reflect"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/golang/mock/gomock"
	"GOFR/models"
	"GOFR/store"
)

var currentTime = time.Now().String()
var patient = models.Patient{
	ID:          1,
	Name:        "ZopSmart",
	Phone:       "+919172681679",
	Discharged:  true,
	CreatedAt:   currentTime,
	UpdatedAt:   currentTime,
	BloodGroup:  "+A",
	Description: "description",
}

func Test_GetByID(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPatientService := store.NewMockPatient(mockCtrl)
	testCases := []struct {
		id            string
		output        models.Patient
		mockCall      *gomock.Call
		expectedError error
		status        int
	}{
		// Success
		{
			id:            "1",
			output:        patient,
			mockCall:      mockPatientService.EXPECT().GetByID(ctx, 1).Return(&patient, nil),
			expectedError: nil,
		},

		// Failure Invalid Id
		{
			id:            "-1",
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}
	p := New(mockPatientService)

	for _, testCase := range testCases {
		_, err := p.GetByID(ctx, testCase.id)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_Create(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPatientService := store.NewMockPatient(mockCtrl)
	testCases := []struct {
		input         models.Patient
		output        models.Patient
		mockCall      *gomock.Call
		expectedError error
		status        int
	}{
		// Success
		{
			input:         patient,
			mockCall:      mockPatientService.EXPECT().Create(ctx, &patient).Return(&patient, nil),
			expectedError: nil,
		},
		// Failure
		{
			input:         patient,
			mockCall:      mockPatientService.EXPECT().Create(ctx, &patient).Return(nil, errors.Error("invalid fileds")),
			expectedError: errors.Error("invalid fileds"),
		},
		// Invalid Id
		{
			input: models.Patient{
				ID:          -1,
				Name:        "ZopSmart",
				Phone:       "+919172681679",
				Discharged:  true,
				CreatedAt:   currentTime,
				UpdatedAt:   currentTime,
				BloodGroup:  "+A",
				Description: "description",
			},
			expectedError: errors.Error("invalid fileds"),
		},
		// Invalid Name
		{
			input: models.Patient{
				ID:          1,
				Name:        "",
				Phone:       "+919172681679",
				Discharged:  true,
				CreatedAt:   currentTime,
				UpdatedAt:   currentTime,
				BloodGroup:  "+A",
				Description: "description",
			},
			expectedError: errors.Error("invalid fileds"),
		},
		// Invalid Phone
		{
			input: models.Patient{
				ID:          1,
				Name:        "ZopSmart",
				Phone:       "+919171679",
				Discharged:  true,
				CreatedAt:   currentTime,
				UpdatedAt:   currentTime,
				BloodGroup:  "+A",
				Description: "description",
			},
			expectedError: errors.Error("invalid fileds"),
		},
	}
	p := New(mockPatientService)

	for _, testCase := range testCases {
		_, err := p.Create(ctx, &testCase.input)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_GetAll(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPatientService := store.NewMockPatient(mockCtrl)
	testCases := []struct {
		output        []models.Patient
		mockCall      *gomock.Call
		expectedError error
		status        int
	}{
		// Success
		{
			output:        []models.Patient{patient},
			mockCall:      mockPatientService.EXPECT().Get(ctx).Return([]*models.Patient{&patient}, nil),
			expectedError: nil,
		},
		// Failure
		{
			mockCall:      mockPatientService.EXPECT().Get(ctx).Return(nil, errors.Error("error")),
			expectedError: errors.Error("error"),
		},
	}

	p := New(mockPatientService)
	for _, testCase := range testCases {
		_, err := p.Get(ctx)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_Update(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPatientService := store.NewMockPatient(mockCtrl)
	testCases := []struct {
		id            string
		output        models.Patient
		mockCall      []*gomock.Call
		expectedError error
		status        int
	}{

		// Success
		{
			id: "1",
			mockCall: []*gomock.Call{mockPatientService.EXPECT().Update(ctx, 1, &patient).Return(&patient, nil),
				mockPatientService.EXPECT().GetByID(ctx, 1).Return(&patient, nil),
			},
			expectedError: nil,
		},
		// Failure Invalid Id
		{
			id:            "-1",
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}
	p := New(mockPatientService)

	for _, testCase := range testCases {
		_, err := p.Update(ctx, testCase.id, &patient)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_Delete(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPatientService := store.NewMockPatient(mockCtrl)
	testCases := []struct {
		id            string
		mockCall      []*gomock.Call
		expectedError error
		status        int
	}{
		// Success
		{
			id: "1",
			mockCall: []*gomock.Call{mockPatientService.EXPECT().Delete(ctx, 1).Return(nil),
				mockPatientService.EXPECT().GetByID(ctx, 1).Return(&patient, nil),
			},
			expectedError: nil,
		},
		// Failure Invalid Id
		{
			id:            "-1",
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}

	p := New(mockPatientService)

	for _, testCase := range testCases {
		err := p.Delete(ctx, testCase.id)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}
