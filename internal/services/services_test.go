package services

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/golang/mock/gomock"
	"github.com/punitj12/patient-app-gofr/internal/models"
	"github.com/punitj12/patient-app-gofr/internal/stores"
	"gopkg.in/guregu/null.v4"
)

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	PatientStoreMock := stores.NewMockPatientStorer(mockCtrl)
	PatientService := New(PatientStoreMock)

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.Context = context.Background()

	tcs := []struct {
		desc    string
		id      int
		patient *models.Patient
		mock    *gomock.Call
		err     error
		result  *models.Patient
	}{
		{
			desc: "success case",
			id:   1,
			patient: &models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			mock: PatientStoreMock.EXPECT().
				Create(ctx, &models.Patient{
					ID:          1,
					Name:        null.NewString("Punit Jain", true),
					Phone:       null.NewString("+916264346285", true),
					BloodGroup:  null.NewString("+B", true),
					Description: null.NewString("Allergic Rhinitis", true),
				}).
				Return(&models.Patient{
					ID:          1,
					Name:        null.NewString("Punit Jain", true),
					Phone:       null.NewString("+916264346285", true),
					BloodGroup:  null.NewString("+B", true),
					Description: null.NewString("Allergic Rhinitis", true),
				}, nil),
			result: &models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
		},
		{
			desc: "validation fail - name",
			id:   2,
			patient: &models.Patient{
				ID:          1,
				Name:        null.StringFrom(""),
				Phone:       null.NewString("+916666612345", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			err: errors.InvalidParam{},
		},
		{
			desc: "validation fail - blood grp",
			id:   2,
			patient: &models.Patient{
				ID:          1,
				Name:        null.StringFrom("Punit"),
				Phone:       null.NewString("+9162as1236285", true),
				BloodGroup:  null.NewString("+Ba", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			err: errors.InvalidParam{},
		},
		{
			desc: "validation fail - phone number",
			id:   3,
			patient: &models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+9162as613346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			err: errors.InvalidParam{},
		},
		{
			desc: "error adding",
			id:   4,
			patient: &models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			mock: PatientStoreMock.EXPECT().
				Create(ctx, &models.Patient{
					ID:          1,
					Name:        null.NewString("Punit Jain", true),
					Phone:       null.NewString("+916264346285", true),
					BloodGroup:  null.NewString("+B", true),
					Description: null.NewString("Allergic Rhinitis", true),
				}).
				Return(&models.Patient{}, errors.Error("error adding patient")),
			err: errors.Error("error adding patient"),
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run("testing add service", func(t *testing.T) {
			{
				res, er := PatientService.Create(ctx, tc.patient)
				if er != nil && (tc.err.Error() != er.Error()) {
					t.Errorf("desc ->>> %v Expected Error : %v Got : %v", tc.desc, tc.err, er)
				}
				if !reflect.DeepEqual(res, tc.result) {
					t.Errorf("desc ->>> %v Expected Result : %v Got : %v", tc.desc, tc.result, res)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	PatientStoreMock := stores.NewMockPatientStorer(mockCtrl)
	PatientService := New(PatientStoreMock)

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.Context = context.Background()

	tcs := []struct {
		desc    string
		id      int
		patient *models.Patient
		mock    *gomock.Call
		err     error
		result  *models.Patient
	}{
		{
			desc: "success",
			id:   1,
			patient: &models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			mock: PatientStoreMock.EXPECT().
				Update(ctx, &models.Patient{
					ID:          1,
					Name:        null.NewString("Punit Jain", true),
					Phone:       null.NewString("+916264346285", true),
					BloodGroup:  null.NewString("+B", true),
					Description: null.NewString("Allergic Rhinitis", true),
				}).
				Return(&models.Patient{
					ID:          1,
					Name:        null.NewString("Punit Jain", true),
					Phone:       null.NewString("+916264346285", true),
					BloodGroup:  null.NewString("+B", true),
					Description: null.NewString("Allergic Rhinitis", true),
				}, nil),
			err: nil,
			result: &models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
		},
		{
			desc: "invalid id",
			id:   2,
			patient: &models.Patient{
				ID:          0,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+9162as64346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			err: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc: "validation fail",
			id:   3,
			patient: &models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+9162as64346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			err: errors.InvalidParam{},
		},
		{
			desc: "error updating",
			id:   4,
			patient: &models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			mock: PatientStoreMock.EXPECT().
				Update(ctx, &models.Patient{
					ID:          1,
					Name:        null.NewString("Punit Jain", true),
					Phone:       null.NewString("+916264346285", true),
					BloodGroup:  null.NewString("+B", true),
					Description: null.NewString("Allergic Rhinitis", true),
				}).
				Return(nil, errors.Error("error updating patient")),
			err: errors.Error("error updating patient"),
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run("testing update service", func(t *testing.T) {
			{
				res, er := PatientService.Update(ctx, tc.patient)
				if er != nil && (tc.err.Error() != er.Error()) {
					t.Errorf("desc ->>> %v Expected Error : %v Got : %v", tc.desc, tc.err, er)
				}
				if !reflect.DeepEqual(res, tc.result) {
					t.Errorf("desc ->>> %v Expected Result : %v Got : %v", tc.desc, tc.result, res)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	PatientStoreMock := stores.NewMockPatientStorer(mockCtrl)
	PatientService := New(PatientStoreMock)

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.Context = context.Background()

	tcs := []struct {
		desc    string
		id      int
		patient models.Patient
		mock    *gomock.Call
		err     error
	}{
		{
			desc: "success case",
			id:   1,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			mock: PatientStoreMock.EXPECT().Delete(ctx, 1),
			err:  nil,
		},
		{
			desc: "invalid id",
			id:   2,
			patient: models.Patient{
				ID:          0,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+9162as64346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			err: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc: "no data found",
			id:   3,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+9162as64346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			mock: PatientStoreMock.EXPECT().Delete(ctx, 1).Return(errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(1)}),
			err:  errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(1)},
		},
		{
			desc: "error deleting",
			id:   4,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			mock: PatientStoreMock.EXPECT().
				Delete(ctx, 1).
				Return(errors.Error("error deleting patient")),
			err: errors.Error("error deleting patient"),
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run("testing update service", func(t *testing.T) {
			{
				er := PatientService.Delete(ctx, tc.patient.ID)
				if er != nil && (tc.err.Error() != er.Error()) {
					t.Errorf("desc ->>> %v Expected Error : %v Got : %v", tc.desc, tc.err, er)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	PatientStoreMock := stores.NewMockPatientStorer(mockCtrl)
	PatientService := New(PatientStoreMock)

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.Context = context.Background()

	tcs := []struct {
		desc    string
		id      int
		patient models.Patient
		mock    *gomock.Call
		err     error
		result  *models.Patient
	}{
		{
			desc: "success case",
			id:   1,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916666612345", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			mock: PatientStoreMock.EXPECT().Get(ctx, 1).
				Return(&models.Patient{
					ID:          1,
					Name:        null.NewString("Punit Jain", true),
					Phone:       null.NewString("+916666612345", true),
					BloodGroup:  null.NewString("+B", true),
					Description: null.NewString("Allergic Rhinitis", true),
				}, nil),
			err: nil,
			result: &models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916666612345", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
		},
		{
			desc: "invalid id",
			id:   2,
			patient: models.Patient{
				ID:          0,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+919as634346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			err: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc: "no data found",
			id:   3,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+915555511111", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			mock: PatientStoreMock.EXPECT().
				Get(ctx, 1).
				Return(nil, errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(1)}),
			err: errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(1)},
		},
		{
			desc: "error deleting",
			id:   4,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+B", true),
				Description: null.NewString("Allergic Rhinitis", true),
			},
			mock: PatientStoreMock.EXPECT().
				Get(ctx, 1).
				Return(&models.Patient{}, errors.Error("error fetching patient")),
			err: errors.Error("error fetching patient"),
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run("testing get service", func(t *testing.T) {
			{
				res, er := PatientService.Get(ctx, tc.patient.ID)
				if er != nil && (tc.err.Error() != er.Error()) {
					t.Errorf("desc ->>> %v Expected Error : %v Got : %v", tc.desc, tc.err, er)
				}
				if !reflect.DeepEqual(res, tc.result) {
					t.Errorf("desc ->>> %v Expected Result : %v Got : %v", tc.desc, tc.result, res)
				}
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	PatientStoreMock := stores.NewMockPatientStorer(mockCtrl)
	PatientService := New(PatientStoreMock)

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.Context = context.Background()

	tcs := []struct {
		desc    string
		id      int
		patient models.Patient
		mock    *gomock.Call
		err     error
		result  []*models.Patient
	}{
		{
			desc: "success case",
			id:   1,
			mock: PatientStoreMock.EXPECT().GetAll(ctx).
				Return([]*models.Patient{
					{
						ID:          1,
						Name:        null.NewString("Punit Jain", true),
						Phone:       null.NewString("+916264346285", true),
						BloodGroup:  null.NewString("+B", true),
						Description: null.NewString("Allergic Rhinitis", true),
					},
					{
						ID:          1,
						Name:        null.NewString("Punit Jain", true),
						Phone:       null.NewString("+916264346285", true),
						BloodGroup:  null.NewString("+B", true),
						Description: null.NewString("Allergic Rhinitis", true),
					}},
					nil),
			err: nil,
			result: []*models.Patient{
				{
					ID:          1,
					Name:        null.NewString("Punit Jain", true),
					Phone:       null.NewString("+916264346285", true),
					BloodGroup:  null.NewString("+B", true),
					Description: null.NewString("Allergic Rhinitis", true),
				},
				{
					ID:          1,
					Name:        null.NewString("Punit Jain", true),
					Phone:       null.NewString("+916264346285", true),
					BloodGroup:  null.NewString("+B", true),
					Description: null.NewString("Allergic Rhinitis", true),
				}},
		},
		{
			desc: "error fetching",
			id:   2,
			mock: PatientStoreMock.EXPECT().GetAll(ctx).Return(nil, errors.Error("error fetching patient data")),
			err:  errors.Error("error fetching patient data"),
		},
		{
			desc: "no data found",
			id:   2,
			mock: PatientStoreMock.EXPECT().GetAll(ctx).Return(nil, errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(2)}),
			err:  errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(2)},
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run("testing get service", func(t *testing.T) {
			{
				res, er := PatientService.GetAll(ctx)
				if er != nil && (tc.err.Error() != er.Error()) {
					t.Errorf("desc ->>> %v Expected Error : %v Got : %v", tc.desc, tc.err, er)
				}
				if !reflect.DeepEqual(res, tc.result) {
					t.Errorf("desc ->>> %v Expected Result : %v Got : %v", tc.desc, tc.result, res)
				}
			}
		})
	}
}
