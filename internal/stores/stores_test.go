package stores

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/punitj12/patient-app-gofr/internal/models"
	"gopkg.in/guregu/null.v4"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Error creating mock : %v", err)
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "phone", "discharged", "bloodGroup", "description", "createdAt", "updatedAt"}).
		AddRow(1, "Punit", "+916666612345", true, "+B", "HeadAche", now, now)

	tcs := []struct {
		desc    string
		id      int
		patient models.Patient
		result  *models.Patient
		mock    interface{}
		err     error
	}{
		{
			desc: "success case",
			id:   1,
			patient: models.Patient{
				ID:          1,
				Name:        null.StringFrom("Punit"),
				Phone:       null.StringFrom("+916666612345"),
				Discharged:  null.BoolFrom(true),
				BloodGroup:  null.StringFrom("+B"),
				Description: null.StringFrom("HeadAche"),
			},
			mock: []interface{}{
				mock.ExpectPrepare("INSERT INTO `patient`(name,phone,discharged,bloodGroup,description) VALUES(?,?,?,?,?);").
					ExpectExec().WithArgs("Punit", "+916666612345", false, "+B", "HeadAche").WillReturnResult(sqlmock.NewResult(1, 1)),

				mock.
					ExpectPrepare(
						"SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from `patient` where id = ? AND deletedAt IS NULL;",
					).
					ExpectQuery().WithArgs(1).WillReturnRows(rows),
			},
			err: nil,
			result: &models.Patient{
				ID:          1,
				Name:        null.StringFrom("Punit"),
				Phone:       null.StringFrom("+916666612345"),
				Discharged:  null.BoolFrom(true),
				BloodGroup:  null.StringFrom("+B"),
				Description: null.StringFrom("HeadAche"),
				CreatedAt:   now, UpdatedAt: now},
		},
		{
			desc: "err executing",
			id:   2,
			patient: models.Patient{
				ID:          1,
				Name:        null.StringFrom("Punit"),
				Phone:       null.StringFrom("+916666612345"),
				Discharged:  null.BoolFrom(true),
				BloodGroup:  null.StringFrom("+B"),
				Description: null.StringFrom("HeadAche"),
			},
			mock: []interface{}{
				mock.ExpectPrepare("INSERT INTO `patient`(name,phone,discharged,bloodGroup,description) VALUES(?,?,?,?,?);").
					ExpectExec().WithArgs("Punit", "+916666612345", true, "+B", "HeadAche").WillReturnError(errors.Error("error executing add query"))},
			err:    errors.Error("error executing add query"),
			result: nil,
		},
		{
			desc: "err preparing",
			id:   3,
			patient: models.Patient{
				ID:          1,
				Name:        null.StringFrom("Punit"),
				Phone:       null.StringFrom("+916666612345"),
				Discharged:  null.BoolFrom(true),
				BloodGroup:  null.StringFrom("+B"),
				Description: null.StringFrom("HeadAche"),
			},
			mock: []interface{}{
				mock.ExpectPrepare("INSERT INTO `patient`(name,phone,discharged,bloodGroup,description) VALUES(?,?,?,?,?);").
					WillReturnError(errors.Error("error preparing add query")),
			},

			err:    errors.Error("error preparing add query"),
			result: nil,
		},
	}

	dbstorer := New()

	for _, tc := range tcs {
		tc := tc

		res, err := dbstorer.Create(ctx, &tc.patient)

		if err != nil && (err.Error() != tc.err.Error()) {
			t.Errorf("desc -----> %v Expected Error : %v Got: %v", tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(tc.result, res) {
			t.Errorf("desc ------> %v Expected Result : %v Got: %v", tc.desc, tc.result, res)
		}
	}
}

func TestGetById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("error creating mock : %v", err)
	}
	defer db.Close()

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "phone", "discharged", "bloodGroup", "description", "createdAt", "updatedAt"}).
		AddRow(1, "Punit Jain", "+916264346285", false, "+b", "Cold", now, now)

	tcs := []struct {
		desc      string
		id        int
		patient   models.Patient
		mockQuery interface{}
		result    *models.Patient
		err       error
	}{
		{
			desc: "success",
			id:   1,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+b", true),
				Description: null.NewString("Cold", true),
			},
			mockQuery: mock.
				ExpectPrepare(
					"SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from `patient` where id = ? AND deletedAt IS NULL;",
				).
				ExpectQuery().WithArgs(1).WillReturnRows(rows),
			err: nil,
			result: &models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				Discharged:  null.NewBool(false, true),
				BloodGroup:  null.NewString("+b", true),
				Description: null.NewString("Cold", true),
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		{
			desc: "err preparing",
			id:   2,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+b", true),
				Description: null.NewString("Cold", true),
			},
			mockQuery: mock.
				ExpectPrepare(
					"SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from `patient` where id = ? AND deletedAt IS NULL;",
				).
				WillReturnError(errors.Error("error preparing select query")),
			err:    errors.Error("error preparing select query"),
			result: nil,
		},
		{
			desc: "no data found",
			id:   4,
			patient: models.Patient{
				ID:          4,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+b", true),
				Description: null.NewString("Cold", true),
			},
			mockQuery: mock.
				ExpectPrepare(
					"SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from `patient` where id = ? AND deletedAt IS NULL;",
				).
				ExpectQuery().WithArgs(4).WillReturnError(sql.ErrNoRows),
			err:    errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(4)},
			result: nil,
		},
		{
			desc: "err querying row",
			id:   3,
			patient: models.Patient{
				ID:          2,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("++91", true),
				BloodGroup:  null.NewString("+b", true),
				Description: null.NewString("Cold", true),
			},
			mockQuery: mock.
				ExpectPrepare(
					"SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from `patient` where id = ? AND deletedAt IS NULL;",
				).
				ExpectQuery().WithArgs(2).WillReturnError(sql.ErrNoRows),
			err:    errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(2)},
			result: nil,
		},
	}

	dbstorer := New()

	for _, tc := range tcs {
		tc := tc
		patient, err := dbstorer.Get(ctx, tc.patient.ID)

		if err != nil && (tc.err.Error() != err.Error()) {
			t.Errorf("desc ---> %v , Expected Error : %v , Got : %v", tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(patient, tc.result) {
			t.Errorf("desc ---> %v, Expected Result : %v , Got : %v", tc.desc, tc.result, patient)
		}
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("error creating mock : %v", err)
	}
	defer db.Close()

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "phone", "discharged", "bloodGroup", "description", "createdAt", "updatedAt"}).
		AddRow(1, "Punit Jain", "+916354346285", false, "+b", "Cold", now, now).
		AddRow(2, "Abhishek Bhandari", "+916666555653", false, "+o", "Cold", now, now)

	errRows := sqlmock.NewRows([]string{"id", "name", "discharged", "bloodGroup", "description", "createdAt", "updatedAt"}).
		AddRow(1, "Punit Jain", false, "+b", "Cold", now, now).
		AddRow(2, "Abhishek Bhandari", false, "+o", "Cold", now, now)

	emptyRows := sqlmock.NewRows([]string{"id", "name", "phone", "discharged", "bloodGroup", "description", "createdAt", "updatedAt"})

	tcs := []struct {
		desc      string
		id        int
		patient   models.Patient
		mockQuery interface{}
		result    []*models.Patient
		err       error
	}{
		{
			desc: "success",
			id:   1,
			mockQuery: mock.
				ExpectPrepare("SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from `patient` where deletedAt IS NULL;").
				ExpectQuery().WillReturnRows(rows),
			err: nil,
			result: []*models.Patient{
				{
					ID:          1,
					Name:        null.NewString("Punit Jain", true),
					Phone:       null.NewString("+916354346285", true),
					Discharged:  null.NewBool(false, true),
					BloodGroup:  null.NewString("+b", true),
					Description: null.NewString("Cold", true),
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          2,
					Name:        null.NewString("Abhishek Bhandari", true),
					Phone:       null.NewString("+916666555653", true),
					Discharged:  null.NewBool(false, true),
					BloodGroup:  null.NewString("+o", true),
					Description: null.NewString("Cold", true),
					CreatedAt:   now,
					UpdatedAt:   now,
				}},
		},
		{
			desc: "no data found",
			id:   1,
			mockQuery: mock.
				ExpectPrepare("SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from `patient` where deletedAt IS NULL;").
				ExpectQuery().WillReturnRows(emptyRows),
			err:    errors.EntityNotFound{Entity: "patient"},
			result: nil,
		},
		{
			desc: "err preparing",
			id:   1,
			mockQuery: mock.
				ExpectPrepare("SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from `patient` where deletedAt IS NULL;").
				WillReturnError(errors.Error("error preparing select query")),
			err:    errors.Error("error preparing select query"),
			result: nil,
		},
		{
			desc: "err querying row",
			id:   1,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+b", true),
				Description: null.NewString("Cold", true),
			},
			mockQuery: mock.
				ExpectPrepare("SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from `patient` where deletedAt IS NULL;").
				ExpectQuery().WillReturnError(errors.Error("error fetching data")),
			err:    errors.Error("error fetching data"),
			result: nil,
		},
		{
			desc: "err scanning row",
			id:   1,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				BloodGroup:  null.NewString("+b", true),
				Description: null.NewString("Cold", true),
			},
			mockQuery: mock.
				ExpectPrepare("SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from `patient` where deletedAt IS NULL;").
				ExpectQuery().WillReturnRows(errRows),
			err:    errors.Error("error scanning data"),
			result: nil,
		},
	}

	dbstorer := New()

	for _, tc := range tcs {
		tc := tc

		patient, err := dbstorer.GetAll(ctx)

		if err != nil && (tc.err.Error() != err.Error()) {
			t.Errorf("desc ---> %v , Expected Error : %v , Got : %v", tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(tc.result, patient) {
			t.Errorf("desc ---> %v, Expected Result : %v, Got : %v", tc.desc, tc.result, patient)
		}
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("error creating mock : %v", err)
	}
	defer db.Close()

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	tcs := []struct {
		desc      string
		id        int
		patient   models.Patient
		mockQuery interface{}
		err       error
	}{
		{
			desc: "success",
			id:   1,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+b", true),
				Description: null.NewString("Cold", true),
			},
			mockQuery: mock.
				ExpectPrepare("UPDATE `patient` SET deletedAt = ? where id = ? AND deletedAt IS NULL").
				ExpectExec().WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1)),
			err: nil,
		},
		{
			desc: "no patient found",
			id:   1,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+b", true),
				Description: null.NewString("Cold", true),
			},
			mockQuery: mock.
				ExpectPrepare("UPDATE `patient` SET deletedAt = ? where id = ? AND deletedAt IS NULL").
				ExpectExec().WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 0)),
			err: errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(1)},
		},
		{
			desc: "err preparing",
			id:   1,
			patient: models.Patient{
				ID:          1,
				Name:        null.NewString("Punit Jain", true),
				Phone:       null.NewString("+916264346285", true),
				BloodGroup:  null.NewString("+b", true),
				Description: null.NewString("Cold", true),
			},
			mockQuery: mock.ExpectPrepare("UPDATE `patient` SET deletedAt = ? where id = ? AND deletedAt IS NULL").
				WillReturnError(errors.Error("error preparing delete query")),
			err: errors.Error("error preparing delete query"),
		},
		{
			desc: "err querying row",
			id:   1,
			mockQuery: mock.ExpectPrepare("UPDATE `patient` SET deletedAt = ? where id = ? AND deletedAt IS NULL").
				ExpectExec().WithArgs(sqlmock.AnyArg(), 1).WillReturnError(errors.Error("error executing delete query")),
			err: errors.Error("error executing delete query"),
		},
	}

	dbstorer := New()

	for _, tc := range tcs {
		tc := tc

		t.Run("testing delete query", func(t *testing.T) {
			err := dbstorer.Delete(ctx, tc.patient.ID)
			if err != nil && (tc.err.Error() != err.Error()) {
				t.Errorf("desc ---> %v , Expected Error : %v \n, Got : %v", tc.desc, tc.err, err)
				return
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Error creating mock : %v", err)
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "phone", "discharged", "bloodGroup", "description", "createdAt", "updatedAt"}).
		AddRow(1, "Punit Jain", "+916666612345", true, "+B", "HeadAche", now, now)

	tcs := []struct {
		desc    string
		id      int
		patient models.Patient
		result  *models.Patient
		mock    interface{}
		err     error
	}{
		{
			desc:    "success case",
			id:      1,
			patient: models.Patient{ID: 1, Name: null.StringFrom("Punit Jain")},
			mock: []interface{}{
				mock.ExpectPrepare("UPDATE `patient` SET name = ? where id = ? AND deletedAt IS NULL;").
					ExpectExec().WithArgs("Punit Jain", 1).WillReturnResult(sqlmock.NewResult(0, 1)),
				mock.ExpectPrepare(`SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from ` +
					"`patient` where id = ? AND deletedAt IS NULL;").
					ExpectQuery().WithArgs(1).WillReturnRows(rows),
			},
			err: nil,
			result: &models.Patient{
				ID:          1,
				Name:        null.StringFrom("Punit Jain"),
				Phone:       null.StringFrom("+916666612345"),
				Discharged:  null.BoolFrom(true),
				BloodGroup:  null.StringFrom("+B"),
				Description: null.StringFrom("HeadAche"),
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		{
			desc:    "rows affected 0",
			id:      1,
			patient: models.Patient{ID: 1, Name: null.StringFrom("Punit Jain")},
			mock: []interface{}{
				mock.ExpectPrepare("UPDATE `patient` SET name = ? where id = ? AND deletedAt IS NULL;").
					ExpectExec().WithArgs("Punit Jain", 1).WillReturnResult(sqlmock.NewResult(0, 0)),
			},
			err: errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(1)},
		},
		{
			desc:    "err executing update",
			id:      1,
			patient: models.Patient{ID: 1, Name: null.StringFrom("Punit Jain")},
			mock: []interface{}{
				mock.ExpectPrepare("UPDATE `patient` SET name = ? where id = ? AND deletedAt IS NULL;").
					ExpectExec().WithArgs("Punit Jain", 1).WillReturnError(errors.Error("error executing update query")),
			},
			err: errors.Error("error executing update query"),
		},
		{
			desc: "err preparing",
			id:   3,
			patient: models.Patient{
				ID:          1,
				Name:        null.StringFrom("Punit"),
				Phone:       null.StringFrom("+916666612345"),
				Discharged:  null.BoolFrom(true),
				BloodGroup:  null.StringFrom("+B"),
				Description: null.StringFrom("HeadAche"),
			},
			mock: []interface{}{
				mock.ExpectPrepare("UPDATE `patient` SET name = ?,phone = ?,discharged = ?," +
					"bloodGroup = ?,description = ? where id= ? AND deletedAt IS NULL;").
					WillReturnError(errors.Error("error preparing update query")),
			},
			err:    errors.Error("error preparing update query"),
			result: nil,
		},
	}

	dbstorer := New()

	for _, tc := range tcs {
		tc := tc
		res, err := dbstorer.Update(ctx, &tc.patient)

		if err != nil && (err.Error() != tc.err.Error()) {
			t.Errorf("desc -----> %v Expected Error : %v Got: %v", tc.desc, tc.err, err)
		}

		if !reflect.DeepEqual(tc.result, res) {
			t.Errorf("desc ------> %v Expected Result : %v Got: %v", tc.desc, tc.result, res)
		}
	}
}
