package patient

import (
	"context"
	"database/sql"
	"log"
	"reflect"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shivanisharma200/patient-management/internal/models"
)

func NewMock() (db *sql.DB, mock sqlmock.Sqlmock, store Patient, ctx *gofr.Context) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
	}

	store = New()
	ctx = gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	return
}

// test for GetById
func Test_GetByID(t *testing.T) {
	db, mock, store, ctx := NewMock()

	q := "SELECT id,name,phone,discharged,created_at,updated_at,blood_group,description from patients WHERE id = ? AND deleted_at IS NULL"

	defer db.Close()

	currentTime := time.Now().String()

	testCases := []struct {
		desc          string
		id            int
		output        *models.Patient
		mockQuery     interface{}
		expectedError error
	}{
		// Success

		{
			desc: "success test case",
			id:   1,
			output: &models.Patient{ID: 1, Name: "ZopSmart", Phone: "+919172681679", Discharged: true,
				CreatedAt: currentTime, UpdatedAt: currentTime, BloodGroup: "+A", Description: "description"},
			mockQuery: mock.
				ExpectQuery(q).
				WithArgs(1).
				WillReturnRows(mock.NewRows([]string{"id", "name", "phone", "discharged", "created_at", "updated_at", "blood_group", "description"}).
					AddRow(1, "ZopSmart", "+919172681679", true, currentTime, currentTime, "+A", "description")),
			expectedError: nil,
		},
		// Failure
		{
			desc: "failure test case",
			id:   1,
			output: &models.Patient{ID: 1, Name: "ZopSmart", Phone: "+919172681679", Discharged: true, CreatedAt: currentTime,
				UpdatedAt: currentTime, BloodGroup: "+A", Description: "description"},
			mockQuery: mock.
				ExpectQuery(q).
				WithArgs(ctx, 1).
				WillReturnError(errors.InvalidParam{Param: []string{"id"}}),
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			_, err := store.GetByID(ctx, testCase.id)

			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
			}
		})
	}
}

// test for Create
func Test_Create(t *testing.T) {
	db, mock, store, ctx := NewMock()

	currentTime := time.Now().String()

	const q = "SELECT id,name,phone,discharged,created_at,updated_at,blood_group,description from patients WHERE id = ? AND deleted_at IS NULL"

	defer db.Close()

	testCases := []struct {
		desc          string
		input         *models.Patient
		output        *models.Patient
		mockQuery     []interface{}
		expectedError error
	}{
		// Success

		{
			desc:  "success test case",
			input: &models.Patient{ID: 1, Name: "ZopSmart", Phone: "+919172681679", Discharged: true, BloodGroup: "+A", Description: "description"},
			output: &models.Patient{ID: 1, Name: "ZopSmart", Phone: "+919172681679", Discharged: true, CreatedAt: currentTime,
				UpdatedAt: currentTime, BloodGroup: "+A", Description: "description"},
			mockQuery: []interface{}{mock.ExpectExec("INSERT INTO patients(name, phone, discharged, blood_group, description) VALUES(?,?,?,?,?)").
				WithArgs("ZopSmart", "+919172681679", true, "+A", "description").
				WillReturnResult(sqlmock.NewResult(1, 1)),
				mock.
					ExpectQuery(q).
					WithArgs(1).
					WillReturnRows(mock.NewRows([]string{"id", "name", "phone", "discharged", "created_at", "updated_at", "blood_group", "description"}).
						AddRow(1, "ZopSmart", "+919172681679", true, currentTime, currentTime, "+A", "description")).WillReturnError(nil),
			},
			expectedError: nil,
		},
		// Failure
		{
			desc:  "failure test case",
			input: &models.Patient{ID: 1, Name: "ZopSmart", Phone: "+919172681679", Discharged: true, BloodGroup: "+A", Description: "description"},
			output: &models.Patient{ID: 1, Name: "ZopSmart", Phone: "+919172681679", Discharged: true, CreatedAt: currentTime,
				UpdatedAt: currentTime, BloodGroup: "+A", Description: "description"},
			mockQuery: []interface{}{mock.ExpectExec("INSERT INTO patients(name, phone, discharged, blood_group, description) VALUES(?,?,?,?,?)").
				WithArgs("ZopSmart", "+919172681679", true, "+A", "description").
				WillReturnError(errors.Error("Failed to create patient")),
				mock.
					ExpectQuery(q).
					WithArgs(1).
					WillReturnError(errors.InvalidParam{Param: []string{"id"}}),
			},
			expectedError: errors.Error("Failed to create patient"),
		},
	}
	for _, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			_, err := store.Create(ctx, testCase.input)

			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
			}
		})
	}
}

// // test for Get
func Test_Get(t *testing.T) {
	db, mock, store, ctx := NewMock()

	defer db.Close()

	currentTime := time.Now().String()
	q := "SELECT id, name, phone, discharged, created_at, updated_at, blood_group, description from patients where deleted_at IS NULL"

	testCases := []struct {
		desc          string
		mockQuery     interface{}
		expectedError error
	}{
		// Success

		{
			desc: "success test case",

			mockQuery: mock.
				ExpectQuery(q).
				WillReturnRows(mock.NewRows([]string{"id", "name", "phone", "discharged", "created_at", "updated_at", "blood_group", "description"}).
					AddRow(1, "ZopSmart", "+919172681679", true, currentTime, currentTime, "+A", "description").
					AddRow(2, "ZopSmart 2", "+919172681679", true, currentTime, currentTime, "+B", "description 2")),
			expectedError: nil,
		},
		// Failure
		{
			desc: "failure test case",
			mockQuery: mock.
				ExpectQuery(q).
				WillReturnError(errors.EntityNotFound{Entity: "Patient"}),
			expectedError: errors.EntityNotFound{Entity: "Patient"},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			_, err := store.Get(ctx)

			if err != nil && err.Error() != testCase.expectedError.Error() {
				t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
			}
		})
	}
}

// test for Update
func Test_Update(t *testing.T) {
	db, mock, store, ctx := NewMock()

	q := "SELECT id,name,phone,discharged,created_at,updated_at,blood_group,description from patients WHERE id = ? AND deleted_at IS NULL"

	defer db.Close()

	currentTime := time.Now().String()

	testCases := []struct {
		desc          string
		id            int
		input         *models.Patient
		output        *models.Patient
		mockQuery     []interface{}
		expectedError error
	}{
		// Success
		{
			desc: "success test case",
			id:   1,
			input: &models.Patient{ID: 1, Name: "ZopSmart", Phone: "+919172681679", Discharged: true, CreatedAt: currentTime,
				UpdatedAt: currentTime, BloodGroup: "+A", Description: "description"},
			output: &models.Patient{ID: 1, Name: "ZopSmart", Phone: "+919172681679", Discharged: true, CreatedAt: currentTime,
				UpdatedAt: currentTime, BloodGroup: "+A", Description: "description"},
			mockQuery: []interface{}{
				mock.ExpectExec("UPDATE patients SET name=?,description=? WHERE id=? AND deleted_at IS NULL").
					WithArgs("ZopSmart", "description", 1).
					WillReturnResult(sqlmock.NewResult(1, 1)),
				mock.
					ExpectQuery(q).
					WithArgs(1).
					WillReturnRows(mock.NewRows([]string{"id", "name", "phone", "discharged", "created_at", "updated_at", "blood_group", "description"}).
						AddRow(1, "ZopSmart", "+919172681679", true, currentTime, currentTime, "+A", "description")),
			},
			expectedError: nil,
		},
		// Failure
		{
			desc: "failure test case",
			id:   1,
			input: &models.Patient{ID: 1, Name: "ZopSmart Updated", Phone: "+919172681679", Discharged: true, CreatedAt: currentTime,
				UpdatedAt: currentTime, BloodGroup: "+A", Description: "description"},
			output: &models.Patient{ID: 1, Name: "ZopSmart Updated", Phone: "+919172681679", Discharged: true, CreatedAt: currentTime,
				UpdatedAt: currentTime, BloodGroup: "+A", Description: "description"},
			mockQuery: []interface{}{mock.ExpectExec("UPDATE patients SET name=?,description=? WHERE id=? AND deleted_at IS NULL").
				WithArgs("ZopSmart Updated", "description", 1).
				WillReturnError(errors.Error("Failed to update patient")),
				mock.
					ExpectQuery(q).
					WithArgs(1).
					WillReturnError(errors.Error("Failed to update patient")),
			},
			expectedError: errors.Error("Failed to update patient"),
		},
	}
	for _, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			_, err := store.Update(ctx, testCase.id, testCase.input)

			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
			}
		})
	}
}

// test for Delete
func Test_Delete(t *testing.T) {
	db, mock, store, ctx := NewMock()

	defer db.Close()

	format := "2006-01-02 15:04:05"

	currentTime := time.Now().Format(format)

	testCases := []struct {
		desc          string
		id            int
		mockQuery     interface{}
		expectedError error
	}{
		// Success

		{
			desc: "success test case",
			id:   1,
			mockQuery: mock.ExpectExec("UPDATE patients SET deleted_at=? WHERE id=? AND deleted_at IS NULL").
				WithArgs(currentTime, 1).
				WillReturnResult(sqlmock.NewResult(1, 1)),
			expectedError: nil,
		},
		// Failure
		{
			desc: "failure test case",
			id:   1,
			mockQuery: mock.ExpectExec("UPDATE patients SET deleted_at=? WHERE id=? AND deleted_at IS NULL").
				WithArgs(currentTime, 1).
				WillReturnError(errors.Error("Failed to delete patient")),
			expectedError: errors.Error("Failed to delete patient"),
		},
	}
	for _, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			err := store.Delete(ctx, testCase.id)

			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
			}
		})
	}
}
