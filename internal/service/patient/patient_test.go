package service

import (
	"context"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/anish-kmr/patient-system/internal/model"
	store "github.com/anish-kmr/patient-system/internal/store/patient"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

var currTime = time.Now()
var patient = model.Patient{
	ID:          1,
	Name:        null.StringFrom("Anish"),
	Phone:       null.StringFrom("+91 9999999999"),
	Discharged:  null.BoolFrom(false),
	BloodGroup:  null.StringFrom("+B"),
	Description: null.StringFrom("Suffering from hypertension"),
	CreatedAt:   null.TimeFrom(currTime),
	UpdatedAt:   null.TimeFrom(currTime),
}

func TestGetPatientById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()

	mockStore := store.NewMockPatient(mockCtrl)

	testcases := []struct {
		description string
		id          int
		out         *model.Patient
		err         error
		mock        *gomock.Call
	}{
		{
			description: "Success Case",
			id:          1,
			out:         &patient,
			err:         nil,
			mock:        mockStore.EXPECT().GetByID(ctx, 1).Return(&patient, nil),
		},
		{
			description: "Id Negative Case",
			id:          -1,
			out:         nil,
			err:         errors.InvalidParam{Param: []string{"id"}},
		},
	}

	patientService := New(mockStore)

	for i, tc := range testcases {
		tc := tc
		i := i

		t.Run(tc.description, func(t *testing.T) {
			out, err := patientService.GetByID(ctx, tc.id)

			assert.Equal(t, out, tc.out, "[Test %v] : Expected :%v Got %v ", i, out, tc.out)
			assert.Equal(t, err, tc.err, "[Test %v] : Expected ERROR :%v Got %v ", i, err, tc.err)
		})
	}
}

func TestGetAllPatient(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := store.NewMockPatient(mockCtrl)

	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()

	queryFilters := map[string]string{"page": "1", "limit": "1"}
	testcases := []struct {
		description string
		filter      map[string]string
		out         []*model.Patient
		err         error
		mock        *gomock.Call
	}{
		{
			description: "Success Case",
			filter:      queryFilters,
			out:         []*model.Patient{&patient},
			err:         nil,
			mock:        mockStore.EXPECT().GetAll(ctx, queryFilters).Return([]*model.Patient{&patient}, nil),
		},
	}

	patientService := New(mockStore)

	for i, tc := range testcases {
		tc := tc
		i := i

		t.Run(tc.description, func(t *testing.T) {
			out, err := patientService.GetAll(ctx, tc.filter)
			assert.Equal(t, out, tc.out, "[Test %v] : Expected :%v Got %v ", i, out, tc.out)
			assert.Equal(t, err, tc.err, "[Test %v] : Expected ERROR :%v Got %v ", i, err, tc.err)
		})
	}
}

func TestCreatePatient(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := store.NewMockPatient(mockCtrl)

	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()

	testcases := []struct {
		description string
		patient     *model.Patient
		out         *model.Patient
		err         error
		mock        *gomock.Call
	}{
		{
			description: "Success Case",
			patient: &model.Patient{
				Name:        null.StringFrom("Anish"),
				Phone:       null.StringFrom("+91 9999999999"),
				BloodGroup:  null.StringFrom("+B"),
				Description: null.StringFrom("Suffering from hypertension"),
			},
			out:  &patient,
			err:  nil,
			mock: mockStore.EXPECT().Create(ctx, gomock.Any()).Return(&patient, nil),
		},
		{
			description: "Name Invalid",
			patient: &model.Patient{
				Name:        null.NewString("", false),
				Phone:       null.StringFrom("+91 9999999999"),
				BloodGroup:  null.StringFrom("+B"),
				Description: null.StringFrom("Suffering from hypertension"),
			},
			out: nil,
			err: errors.InvalidParam{Param: []string{"name"}},
		},
		{
			description: "Phone Invalid",
			patient: &model.Patient{
				Name:        null.StringFrom("Anish"),
				Phone:       null.StringFrom("+91 99999999a"),
				BloodGroup:  null.StringFrom("+B"),
				Description: null.StringFrom("Suffering from hypertension"),
			},
			out: nil,
			err: errors.InvalidParam{Param: []string{"phone"}},
		},
		{
			description: "Blood Group Invalid",
			patient: &model.Patient{
				Name:        null.StringFrom("Anish"),
				Phone:       null.StringFrom("+91 9999999999"),
				BloodGroup:  null.StringFrom(""),
				Description: null.StringFrom("Suffering from hypertension"),
			},
			out: nil,
			err: errors.InvalidParam{Param: []string{"bloodgroup"}},
		},
	}

	patientService := New(mockStore)

	for i, tc := range testcases {
		tc := tc
		i := i

		t.Run(tc.description, func(t *testing.T) {
			out, err := patientService.Create(ctx, tc.patient)
			assert.Equal(t, out, tc.out, "[Test %v] : Expected :%v Got %v ", i, out, tc.out)
			assert.Equal(t, err, tc.err, "[Test %v] : Expected ERROR :%v Got %v ", i, err, tc.err)
		})
	}
}

func TestUpdatePatient(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := store.NewMockPatient(mockCtrl)
	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	testcases := []struct {
		description string
		id          int
		patient     *model.Patient
		out         *model.Patient
		err         error
		mock        interface{}
	}{
		{
			description: "Success Case",
			id:          1,
			patient:     &patient,
			out:         &patient,
			err:         nil,
			mock: []interface{}{
				mockStore.EXPECT().GetByID(ctx, 1).Return(&patient, nil),
				mockStore.EXPECT().Update(ctx, 1, gomock.Any()).Return(&patient, nil),
			},
		},
		{
			description: "No Patient",
			id:          1,
			patient:     &patient,
			out:         nil,
			err:         errors.EntityNotFound{Entity: "Patient", ID: "1"},
			mock: []interface{}{
				mockStore.EXPECT().GetByID(ctx, 1).Return(nil, errors.EntityNotFound{Entity: "Patient", ID: "1"}),
			},
		},
		{
			description: "Name Invalid",
			patient: &model.Patient{
				ID:          1,
				Name:        null.StringFrom(""),
				Phone:       null.StringFrom("+91 9999999999"),
				Discharged:  null.BoolFrom(false),
				BloodGroup:  null.StringFrom("+B"),
				Description: null.StringFrom("Suffering from hypertension"),
				CreatedAt:   null.TimeFrom(currTime),
				UpdatedAt:   null.TimeFrom(currTime),
			},
			out: nil,
			err: errors.InvalidParam{Param: []string{"name"}},
		},
		{
			description: "Phone Invalid",
			patient: &model.Patient{
				ID:          1,
				Name:        null.StringFrom("Anish"),
				Phone:       null.StringFrom("+91 9a9a9a9a9a"),
				Discharged:  null.BoolFrom(false),
				BloodGroup:  null.StringFrom("+B"),
				Description: null.StringFrom("Suffering from hypertension"),
				CreatedAt:   null.TimeFrom(currTime),
				UpdatedAt:   null.TimeFrom(currTime),
			},
			out: nil,
			err: errors.InvalidParam{Param: []string{"phone"}},
		},
		{
			description: "Blood Group Invalid",
			patient: &model.Patient{
				Name:        null.StringFrom("Anish"),
				Phone:       null.StringFrom("+91 9999999999"),
				BloodGroup:  null.StringFrom(""),
				Description: null.StringFrom("Suffering from hypertension"),
			},
			out: nil,
			err: errors.InvalidParam{Param: []string{"bloodgroup"}},
		},
	}

	patientService := New(mockStore)

	for i, tc := range testcases {
		tc := tc
		i := i

		t.Run(tc.description, func(t *testing.T) {
			out, err := patientService.Update(ctx, tc.id, tc.patient)
			assert.Equal(t, out, tc.out, "[Test %v] : Expected :%v Got %v ", i, out, tc.out)
			assert.Equal(t, err, tc.err, "[Test %v] : Expected ERROR :%v Got %v ", i, err, tc.err)
		})
	}
}

func TestDeletePatient(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := store.NewMockPatient(mockCtrl)
	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	testcases := []struct {
		description string
		id          int
		err         error
		mock        interface{}
	}{
		{
			description: "Success Case",
			id:          1,
			err:         nil,

			mock: []interface{}{
				mockStore.EXPECT().GetByID(ctx, 1).Return(&patient, nil),
				mockStore.EXPECT().Delete(ctx, 1).Return(nil),
			},
		},
		{
			description: "No Patient",
			id:          1,
			err:         errors.EntityNotFound{Entity: "Patient", ID: "1"},

			mock: []interface{}{
				mockStore.EXPECT().GetByID(ctx, 1).Return(nil, errors.EntityNotFound{Entity: "Patient", ID: "1"}),
			},
		},
		{
			description: "Id Negative Case",
			id:          -1,
			err:         errors.InvalidParam{Param: []string{"id"}},
		},
	}

	patientService := New(mockStore)

	for i, tc := range testcases {
		tc := tc
		i := i

		t.Run(tc.description, func(t *testing.T) {
			err := patientService.Delete(ctx, tc.id)
			assert.Equal(t, err, tc.err, "[Test %v] : Expected ERROR :%v Got %v ", i, err, tc.err)
		})
	}
}
