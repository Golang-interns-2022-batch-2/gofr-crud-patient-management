package store

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/anish-kmr/patient-system/internal/model"
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

func newMock() (db *sql.DB, mock sqlmock.Sqlmock, store *PatientStore, ctx *gofr.Context) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		fmt.Println(err)
	}

	store = New()
	ctx = gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	return
}
func TestGetById(t *testing.T) {
	db, mock, patientStore, ctx := newMock()

	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name", "phone", "discharged", "bloodgroup", "description", "createdAt", "updatedAt"})

	testcases := []struct {
		description string
		id          int
		out         *model.Patient
		err         error
		mockQ       interface{}
	}{
		{
			description: "Success Case",
			id:          1,
			out:         &patient,
			err:         nil,
			mockQ: mock.ExpectPrepare(
				"SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE id=? AND deletedAt IS NULL",
			).
				ExpectQuery().
				WithArgs(1).
				WillReturnRows(
					row.AddRow(1, "Anish", "+91 9999999999", false, "+B", "Suffering from hypertension", currTime, currTime),
				),
		},
		{
			description: "Fail Prepare",
			id:          1,
			out:         nil,
			err:         errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare(
				"SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE id=? AND deletedAt IS NULL",
			).
				WillReturnError(
					errors.Error("internal server error"),
				),
		},
		{
			description: "Fail Query",
			id:          1,
			out:         nil,
			err:         errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare(
				"SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE id=? AND deletedAt IS NULL",
			).
				ExpectQuery().
				WithArgs(1).
				WillReturnError(
					errors.Error("internal server error"),
				),
		},
		{
			description: "No Rows Found",
			id:          1,
			out:         nil,
			err:         errors.EntityNotFound{Entity: "Patient", ID: "1"},
			mockQ: mock.ExpectPrepare(
				"SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE id=? AND deletedAt IS NULL",
			).
				ExpectQuery().
				WithArgs(1).
				WillReturnError(
					sql.ErrNoRows,
				),
		},
	}

	for i, tc := range testcases {
		tc := tc
		i := i

		t.Run(tc.description, func(t *testing.T) {
			out, err := patientStore.GetByID(ctx, tc.id)
			assert.Equal(t, out, tc.out, "[Test %v] : Expected :%v Got %v ", i, out, tc.out)
			assert.Equal(t, err, tc.err, "[Test %v] : Expected ERROR :%v Got %v ", i, err, tc.err)
		})
	}
}

func TestGetAll(t *testing.T) {
	db, mock, patientStore, ctx := newMock()

	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name", "phone", "discharged", "bloodgroup", "description", "createdAt", "updatedAt"})
	testcases := []struct {
		description string
		page        int
		limit       int
		filters     map[string]string
		out         []*model.Patient
		err         error
		mockQ       interface{}
	}{
		{
			description: "Success Case",
			filters:     map[string]string{"page": "0", "limit": "5"},
			out:         []*model.Patient{&patient},
			err:         nil,
			mockQ: mock.ExpectPrepare(
				"SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE deletedAt IS NULL LIMIT ?,?",
			).
				ExpectQuery().
				WithArgs(0, 5).
				WillReturnRows(
					row.AddRow(1, "Anish", "+91 9999999999", false, "+B", "Suffering from hypertension", currTime, currTime),
				),
		},
		{
			description: "Scan Error",
			filters:     map[string]string{"page": "0", "limit": "5"},
			out:         nil,
			err:         errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare(
				"SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE deletedAt IS NULL LIMIT ?,?",
			).
				ExpectQuery().
				WithArgs(0, 5).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "phone", "discharged", "bloodgroup", "description"}).
						AddRow(1, "Anish", "+91 9999999999", false, "+B", "Suffering from hypertension"),
				),
		},
		{
			description: "Fail Query",
			filters:     map[string]string{},
			out:         nil,
			err:         errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare(
				"SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE deletedAt IS NULL LIMIT ?,?",
			).
				ExpectQuery().
				WithArgs(0, 5).
				WillReturnError(
					errors.Error("internal server error"),
				),
		},
		{
			description: "Fail Prepare",
			filters:     map[string]string{"page": "0", "limit": "1", "discharged": "0"},
			out:         nil,
			err:         errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare(
				"SELECT " +
					"id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt" +
					"FROM Patient " +
					"WHERE deletedAt IS NULL AND discharged=? LIMIT ?,?",
			).
				WillReturnError(
					errors.Error("internal server error"),
				),
		},
	}

	for i, tc := range testcases {
		tc := tc
		i := i

		t.Run(tc.description, func(t *testing.T) {
			out, err := patientStore.GetAll(ctx, tc.filters)
			assert.Equal(t, out, tc.out, "[Test %v] : Expected :%v Got %v ", i, out, tc.out)
			assert.Equal(t, err, tc.err, "[Test %v] : Expected ERROR :%v Got %v ", i, err, tc.err)
		})
	}
}

func TestCreate(t *testing.T) {
	db, mock, patientStore, ctx := newMock()

	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name", "phone", "discharged", "bloodgroup", "description", "createdAt", "updatedAt"})
	testcases := []struct {
		description string
		patient     *model.Patient
		out         *model.Patient
		err         error
		mockQ       interface{}
	}{
		{
			description: "Success Case",
			patient: &model.Patient{
				Name:        patient.Name,
				Phone:       patient.Phone,
				BloodGroup:  patient.BloodGroup,
				Description: patient.Description,
			},
			out: &patient,
			err: nil,
			mockQ: []interface{}{
				mock.ExpectPrepare("INSERT INTO Patient(name,phone,bloodgroup,description) Values(?,?,?,?)").
					ExpectExec().
					WithArgs(patient.Name, patient.Phone, patient.BloodGroup, patient.Description).
					WillReturnResult(
						sqlmock.NewResult(1, 1),
					),
				mock.ExpectPrepare(
					"SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE id=? AND deletedAt IS NULL",
				).
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(
						row.AddRow(patient.ID,
							patient.Name,
							patient.Phone,
							patient.Discharged,
							patient.BloodGroup,
							patient.Description,
							patient.CreatedAt,
							patient.UpdatedAt,
						),
					),
			},
		},
		{
			description: "Fail Prepare",
			out:         nil,
			err:         errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare("INSERT INTO Patient(name,phone,bloodgroup,description) Values(?,?,?,?)").
				WillReturnError(errors.Error("internal server error")),
		},
		{
			description: "Fail Execute",
			patient: &model.Patient{
				Name:        patient.Name,
				Phone:       patient.Phone,
				BloodGroup:  patient.BloodGroup,
				Description: patient.Description,
			},
			out: nil,
			err: errors.Error("internal server error"),
			mockQ: mock.
				ExpectPrepare("INSERT INTO Patient(name,phone,bloodgroup,description) Values(?,?,?,?)").
				ExpectExec().
				WithArgs(patient.Name, patient.Phone, patient.BloodGroup, patient.Description).
				WillReturnError(errors.Error("internal server error")),
		},
	}

	for i, tc := range testcases {
		tc := tc
		i := i

		t.Run(tc.description, func(t *testing.T) {
			out, err := patientStore.Create(ctx, tc.patient)

			assert.Equal(t, out, tc.out, "[Test %v] : Expected :%v Got %v ", i, out, tc.out)
			assert.Equal(t, err, tc.err, "[Test %v] : Expected ERROR :%v Got %v ", i, err, tc.err)
		})
	}
}
func TestUpdate(t *testing.T) {
	db, mock, patientStore, ctx := newMock()

	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name", "phone", "discharged", "bloodgroup", "description", "createdAt", "updatedAt"})
	testcases := []struct {
		description string
		id          int
		patient     *model.Patient
		out         *model.Patient
		err         error
		mockQ       interface{}
	}{
		{
			description: "Success Case - Partial Fields",
			id:          1,
			patient: &model.Patient{
				Name: null.StringFrom("Anish"),
			},
			out: &patient,
			err: nil,
			mockQ: []interface{}{
				mock.ExpectPrepare("UPDATE Patient Set name=? WHERE id=? AND deletedAt IS NULL").
					ExpectExec().
					WithArgs(
						patient.Name,
						1,
					).
					WillReturnResult(
						sqlmock.NewResult(1, 1),
					),
				mock.ExpectPrepare(
					"SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE id=? AND deletedAt IS NULL",
				).
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(
						row.AddRow(
							patient.ID,
							patient.Name,
							patient.Phone,
							patient.Discharged,
							patient.BloodGroup,
							patient.Description,
							patient.CreatedAt,
							patient.UpdatedAt,
						),
					),
			},
		},

		{
			description: "Success Case",
			id:          1,
			patient:     &patient,
			out:         &patient,
			err:         nil,
			mockQ: []interface{}{
				mock.ExpectPrepare("UPDATE Patient Set name=? ,phone=? ,bloodgroup=? ,discharged=? ,description=? WHERE id=? AND deletedAt IS NULL").
					ExpectExec().
					WithArgs(
						patient.Name,
						patient.Phone,
						patient.BloodGroup,
						patient.Discharged,
						patient.Description,
						1,
					).
					WillReturnResult(
						sqlmock.NewResult(1, 1),
					),
				mock.ExpectPrepare(
					"SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE id=? AND deletedAt IS NULL",
				).
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(
						row.AddRow(
							patient.ID,
							patient.Name,
							patient.Phone,
							patient.Discharged,
							patient.BloodGroup,
							patient.Description,
							patient.CreatedAt,
							patient.UpdatedAt,
						),
					),
			},
		},
		{
			description: "No Fields",
			patient: &model.Patient{
				Name:        null.NewString("", false),
				Phone:       null.NewString("", false),
				Description: null.NewString("", false),
				Discharged:  null.NewBool(false, false),
				BloodGroup:  null.NewString("", false),
			},
			out: nil,
			err: errors.InvalidParam{Param: []string{"body"}},
		},
		{
			description: "No Rows Affected Case",
			id:          1,
			patient:     &patient,
			out:         nil,
			err:         errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare(
				"UPDATE Patient Set name=? ,phone=? ,bloodgroup=? ,discharged=? ,description=? WHERE id=? AND deletedAt IS NULL",
			).
				ExpectExec().
				WithArgs(
					patient.Name,
					patient.Phone,
					patient.BloodGroup,
					patient.Discharged,
					patient.Description,
					1,
				).
				WillReturnResult(
					sqlmock.NewResult(0, 0),
				),
		},
		{
			description: "Fail Exec",
			patient:     &patient,
			out:         nil,
			err:         errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare(
				"UPDATE Patient Set name=? ,phone=? ,bloodgroup=? ,discharged=? ,description=? WHERE id=? AND deletedAt IS NULL",
			).
				ExpectExec().
				WithArgs(
					patient.Name,
					patient.Phone,
					patient.BloodGroup,
					patient.Discharged,
					patient.Description,
					1,
				).
				WillReturnError(errors.Error("internal server error")),
		},
		{
			description: "Fail Prepare",
			patient:     &patient,
			out:         nil,
			err:         errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare(
				"UPDATE Patient SET name=?, phone=?, discharged=?, bloodgroup=?, description=? WHERE id=? AND deletedAt is NULL",
			).
				WillReturnError(errors.Error("internal server error")),
		},
	}

	for i, tc := range testcases {
		tc := tc
		i := i

		t.Run(tc.description, func(t *testing.T) {
			out, err := patientStore.Update(ctx, tc.id, tc.patient)

			assert.Equal(t, out, tc.out, "[Test %v] : Expected :%v Got %v ", i, out, tc.out)
			assert.Equal(t, err, tc.err, "[Test %v] : Expected ERROR :%v Got %v ", i, err, tc.err)
		})
	}
}

func TestDelete(t *testing.T) {
	db, mock, patientStore, ctx := newMock()

	defer db.Close()

	testcases := []struct {
		description string
		id          int
		err         error
		mockQ       interface{}
	}{
		{
			description: "Success Case",
			id:          1,
			err:         nil,
			mockQ: mock.ExpectPrepare("UPDATE Patient SET deletedAt=? WHERE id=? AND deletedAt is NULL").
				ExpectExec().
				WithArgs(sqlmock.AnyArg(), 1).
				WillReturnResult(
					sqlmock.NewResult(1, 1),
				),
		},

		{
			description: "Fail Prepare",
			id:          1,
			err:         errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare("UPDATE Patient SET deletedAt=? WHERE id=? AND deletedAt is NULL").
				WillReturnError(errors.Error("internal server error")),
		},
		{
			description: "Fail Execute",
			id:          1,
			err:         errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare("UPDATE Patient SET deletedAt=? WHERE id=? AND deletedAt is NULL").
				ExpectExec().
				WithArgs(sqlmock.AnyArg(), 1).
				WillReturnError(errors.Error("internal server error")),
		},
	}

	for i, tc := range testcases {
		tc := tc
		i := i

		t.Run(tc.description, func(t *testing.T) {
			err := patientStore.Delete(ctx, tc.id)
			assert.Equal(t, err, tc.err, "[Test %v] : Expected err :%v Got %v ", i, err, tc.err)
		})
	}
}
