package patient

import (
	"context"
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aakanksha/updated-patient-management-system/internal/models"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestGetById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	tests := []struct {
		id          int
		input       int
		output      *models.Patient
		mockQuery   interface{}
		expectError error
	}{

		{
			id:     3,
			input:  3,
			output: &models.Patient{Id: 3},

			mockQuery:   mock.ExpectQuery("select id,name,phone,discharge,createdat,udatedat,bloodgroup,description from patient where deletedat IS NULL and id=?").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone", "discharge", "createdat", "udatedat", "bloodgroup", "description"}).AddRow(1, "ak", "123", true, "2006-02-01 15:04:05", "2006-02-01 15:04:05", "a+", "good")),
			expectError: errors.New("id not found"),
		},
		{
			id:          4,
			input:       4,
			output:      &models.Patient{},
			mockQuery:   mock.ExpectExec("select id,name,phone,discharge,createdat,udatedat,bloodgroup,description from patient where deletedat IS NULL and id=?").WithArgs(4).WillReturnError(errors.New("id not found")),
			expectError: errors.New("id not found"),
		},
		{
			id:     3,
			input:  3,
			output: &models.Patient{Id: 3, Name: "aakanksha3", Phone: "123", Discharge: true, BloodGroup: "A+", Description: "abc"},

			mockQuery:   mock.ExpectQuery("select id,name,phone,discharge,createdat,udatedat,bloodgroup,description from patient where deletedat IS NULL and id=?").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone", "discharge", "createdat", "udatedat", "bloodgroup", "description"}).AddRow(1, "ak", "123", true, "02-01-2006 15:04:05", "2006-02-01 15:04:05", "a+", "good")).WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
		},
		{
			id:     3,
			input:  3,
			output: &models.Patient{Id: 3, Name: "aakanksha3", Phone: "123", Discharge: true, BloodGroup: "A+", Description: "abc"},

			mockQuery:   mock.ExpectQuery("select id,name,phone,discharge,createdat,udatedat,bloodgroup,description from patient where deletedat IS NULL and id=?").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone", "discharge", "createdat", "udatedat", "bloodgroup", "description"}).AddRow(1, "ak", "123", true, "2006-31-31 15:04:05", "02-01-2006 15:04:05", "a+", "good")).WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
		},
	}

	for _, testCase := range tests {
		t.Run("", func(t *testing.T) {

			a := New()
			ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
			ctx.Context = context.Background()
			_, err := a.GetByID(ctx, testCase.id)
			if err != nil {
				fmt.Println(err)
			}
			if !reflect.DeepEqual(err, testCase.expectError) {
				t.Errorf("expected error :%v, got :%v ", testCase.expectError, err)
			}

		})
		return
	}
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testcases := []struct {
		id          int
		input       *models.Patient
		output      *models.Patient
		mockQuery   interface{}
		expectError error
	}{
		{
			id:     3,
			input:  &models.Patient{Name: "aakanksha3", Phone: "123", Discharge: true, BloodGroup: "A+", Description: "abc"},
			output: &models.Patient{Name: "aakanksha3", Phone: "123", Discharge: true, BloodGroup: "A+", Description: "abc"},

			mockQuery:   mock.ExpectExec("insert into patient (name,phone,discharge,bloodgroup,description) values (?, ?, ?, ?, ?)").WithArgs("aakanksha3", "123", true, "A+", "abc").WillReturnResult(sqlmock.NewResult(3, 1)),
			expectError: nil,
		},
		{
			id:          4,
			input:       &models.Patient{Name: "aakanksha3", Phone: "123", Discharge: true, BloodGroup: "A+", Description: "abc"},
			output:      &models.Patient{Name: "aakanksha3", Phone: "123", Discharge: true, BloodGroup: "A+", Description: "abc"},
			mockQuery:   mock.ExpectExec("insert into patient (name,phone,discharge,bloodgroup,description) values (?, ?, ?, ?, ?)").WithArgs("aakanksha3", "123", true, "A+", "abc").WillReturnResult(sqlmock.NewResult(3, 1)),
			expectError: nil,
		},
		{
			id:          4,
			input:       &models.Patient{},
			output:      &models.Patient{},
			mockQuery:   mock.ExpectExec("insert into patient (name,phone,discharge,bloodgroup,description) values (?, ?, ?, ?, ?)").WithArgs("aakanksha3", "123", true, "A+").WillReturnError(errors.New("error in executing insert")),
			expectError: errors.New("error in executing insert"),
		},
	}

	for _, testCase := range testcases {
		t.Run("", func(t *testing.T) {

			a := New()
			ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
			ctx.Context = context.Background()
			_, err := a.Insert(ctx, testCase.input)
			if err != nil {
				fmt.Println(err)
			}

			if err != nil && err.Error() != testCase.expectError.Error() {
				t.Errorf("expected error :%v, got :%v ", testCase.expectError, err)
			}

		})
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		id          int
		input       *models.Patient
		output      *models.Patient
		mockQuery   interface{}
		expectError error
	}{
		{
			id:     3,
			input:  &models.Patient{Id: 3, Name: "aakanksha3", Phone: "123", Discharge: true, BloodGroup: "A+", Description: "abc"},
			output: &models.Patient{Id: 3, Name: "aakanksha3", Phone: "123", Discharge: true, BloodGroup: "A+", Description: "abc"},

			mockQuery:   mock.ExpectExec("update patient SET name = ?, phone=?, discharge=?,bloodgroup=?,description=? where deletedat IS NULL and id=?").WithArgs("aakanksha3", "123", true, "A+", "abc", int64(3)).WillReturnResult(sqlmock.NewResult(3, 1)),
			expectError: nil,
		},
		{
			id:          4,
			input:       &models.Patient{Name: "", Phone: "", Discharge: true, BloodGroup: "", Description: ""},
			output:      &models.Patient{Name: "", Phone: "", Discharge: true, BloodGroup: "", Description: ""},
			mockQuery:   mock.ExpectExec("update patient SET name = ?, phone=?, discharge=?,bloodgroup=?,description=? where deletedat IS NULL and id=?").WithArgs("", "", true, "", "", int64(3)).WillReturnError(errors.New("update failed")),
			expectError: errors.New("update failed"),
		},
	}

	for _, testCase := range tests {
		t.Run("", func(t *testing.T) {

			a := New()
			ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
			ctx.Context = context.Background()
			_, err := a.Update(ctx, testCase.input)
			if err != nil {
				fmt.Println(err)
			}

			if err != nil && err.Error() != testCase.expectError.Error() {
				t.Errorf("expected error :%v, got :%v ", testCase.expectError, err)
			}

		})
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	//updaterr := errors.New("update failed")
	format := "2006-01-02 15:04:05"
	tests := []struct {
		id          int
		input       int
		mockQuery   interface{}
		expectError error
	}{
		{
			id:          1,
			input:       1,
			mockQuery:   mock.ExpectExec("UPDATE patient SET deletedat=? WHERE id=? AND deletedat IS NULL").WithArgs(time.Now().Format(format), 1).WillReturnResult(sqlmock.NewResult(1, 1)),
			expectError: nil,
		},
		{
			id:          4,
			input:       4,
			mockQuery:   mock.ExpectExec("UPDATE patient SET deletedat=? WHERE id=? AND deletedat IS NULL").WithArgs(time.Now().Format(format), 4).WillReturnError(errors.New("error of delete")),
			expectError: errors.New("error of delete"),
		},
	}

	for _, testCase := range tests {
		t.Run("", func(t *testing.T) {

			a := New()
			ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
			ctx.Context = context.Background()
			err := a.Delete(ctx, testCase.input)
			if err != nil {
				fmt.Println(err)
			}
			if err != nil && err.Error() != testCase.expectError.Error() {
				t.Errorf("expected error :%v, got :%v ", testCase.expectError, err)
			}

		})
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	createat := time.Now()
	updateat := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "phone", "discharge", "createdAt", "updatedAt", "bloodGroup", "description"}).
		AddRow(1, "P", "+916354346285", false, createat, updateat, "+b", "Cold").
		AddRow(2, "a", "+916666555653", false, createat, updateat, "+o", "Cold")

	tests := []struct {
		output      []*models.Patient
		mockQuery   interface{}
		expectError error
	}{
		{
			output:      []*models.Patient{{Id: 1, Name: "aakanksha3", Phone: "123", Discharge: true, BloodGroup: "A+", Description: "abc"}},
			mockQuery:   mock.ExpectQuery("select id,name,phone,discharge,createdat,udatedat,bloodgroup,description from patient where deletedat IS NULL;").WillReturnRows(rows),
			expectError: nil,
		},
		//{
		//	output:      []*models.Patient{{Id: 3, Name: "aakanksha3", Phone: "123", Discharge: true, BloodGroup: "A+", Description: "abc"}},
		//	mockQuery:   mock.ExpectQuery("select id,name,phone,discharge,createdat,udatedat,bloodgroup,description from patient where deletedat IS NULL;").WillReturnError("not passesd correct data"),
		//	expectError: errors.New("not passesd correct data"),
		//},
	}

	for _, testCase := range tests {
		t.Run("", func(t *testing.T) {

			a := New()
			ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
			ctx.Context = context.Background()
			_, err := a.GetAll(ctx)
			if err != nil {
				fmt.Println(err)
			}

			if err != nil && err.Error() != testCase.expectError.Error() {
				t.Errorf("expected error :%v, got :%v ", testCase.expectError, err)
			}

		})
	}
}
