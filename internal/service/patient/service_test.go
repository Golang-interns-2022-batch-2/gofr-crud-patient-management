package patient

import (
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"errors"
	"github.com/aakanksha/updated-patient-management-system/internal/models"
	"github.com/aakanksha/updated-patient-management-system/internal/stores"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

var patient = models.Patient{
	Id:          1,
	Name:        "aakanksha",
	Phone:       "8083860404",
	Discharge:   true,
	BloodGroup:  "A+",
	Description: "good",
}

func TestService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbhandler := stores.NewMockStoreInterface(ctrl)
	serv := New(dbhandler)
	app := gofr.New()
	testcases := []struct {
		desc   string
		inp    int
		exp    *models.Patient
		expErr error
		mock   []*gomock.Call
	}{
		{
			"testcase-1",
			1,
			&patient,
			nil,
			[]*gomock.Call{
				dbhandler.EXPECT().GetByID(gomock.Any(), 1).Return(&patient, nil),
			},
		},
		{
			"testcase-2",
			1,
			&models.Patient{},
			errors.New("error found"),
			[]*gomock.Call{
				dbhandler.EXPECT().GetByID(gomock.Any(), 1).Return(&models.Patient{}, errors.New("error found")),
			},
		},
	}
	for _, tcs := range testcases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": strconv.Itoa(tcs.inp),
		})
		_, err := serv.GetByID(ctx, tcs.inp)

		if !reflect.DeepEqual(err, tcs.expErr) {
			t.Errorf("%v, expected %v, got %v", tcs.desc, tcs.expErr, err)
		}

	}
}

func Test_DeleteById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	dbhandler := stores.NewMockStoreInterface(mockCtrl)
	app := gofr.New()
	testCases := []struct {
		id            int
		mockCall      *gomock.Call
		expectedError error
		status        int
	}{

		{
			id:            1,
			mockCall:      dbhandler.EXPECT().Delete(gomock.Any(), 1).Return(nil),
			expectedError: nil,
		},

		{
			id:            -1,
			mockCall:      dbhandler.EXPECT().Delete(gomock.Any(), -1).Return(errors.New("unable to delete user")),
			expectedError: errors.New("unable to delete user"),
		},
	}
	p := New(dbhandler)

	for _, testCase := range testCases {

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": strconv.Itoa(testCase.id),
		})
		err := p.Delete(ctx, testCase.id)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func TestUserService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbhandler := stores.NewMockStoreInterface(ctrl)
	serv := New(dbhandler)
	app := gofr.New()
	testcases := []struct {
		desc   string
		exp    []*models.Patient
		expErr error
		mock   []*gomock.Call
	}{
		{
			"testcase-1",
			[]*models.Patient{&patient},
			nil,
			[]*gomock.Call{
				dbhandler.EXPECT().GetAll(gomock.Any()).Return([]*models.Patient{&patient}, nil),
			},
		},
		{
			"testcase-2",
			[]*models.Patient{},
			sql.ErrNoRows,
			[]*gomock.Call{
				dbhandler.EXPECT().GetAll(gomock.Any()).Return([]*models.Patient{}, sql.ErrNoRows),
			},
		},
	}
	for _, tcs := range testcases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)
		_, err := serv.GetAll(ctx)

		if !reflect.DeepEqual(err, tcs.expErr) {
			t.Errorf("%v, expected %v, got %v", tcs.desc, tcs.expErr, err)
		}

	}
}

func TestInsertpatient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbhandler := stores.NewMockStoreInterface(ctrl)
	shandler := New(dbhandler)
	app := gofr.New()
	testcases := []struct {
		desc   string
		inp    models.Patient
		exp    *models.Patient
		expErr error
		mock   []*gomock.Call
	}{
		{
			"testcase-1",
			patient,
			&patient,
			nil,
			[]*gomock.Call{
				dbhandler.EXPECT().Insert(gomock.Any(), &patient).Return(&patient, nil),
				dbhandler.EXPECT().GetByID(gomock.Any(), 1).Return(&patient, nil),
			},
		},
	}
	for _, tcs := range testcases {

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := shandler.Insert(ctx, &tcs.inp)

		if !reflect.DeepEqual(err, tcs.expErr) {
			t.Errorf("%v, expected %v, got %v", tcs.desc, tcs.expErr, err)
		}

	}
}

func Test_UpdateById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	dbhandler := stores.NewMockStoreInterface(mockCtrl)
	serv := New(dbhandler)
	app := gofr.New()
	testCases := []struct {
		id            int
		input         *models.Patient
		output        models.Patient
		mockCall      []*gomock.Call
		expectedError error
		status        int
	}{
		// Success
		{
			id:     1,
			input:  &patient,
			output: patient,
			mockCall: []*gomock.Call{
				dbhandler.EXPECT().Update(gomock.Any(), &patient).Return(&patient, nil),
				dbhandler.EXPECT().GetByID(gomock.Any(), 1).Return(&patient, nil),
			},
			expectedError: nil,
		},
		// Failure
		{
			id:    2,
			input: &patient,
			mockCall: []*gomock.Call{

				dbhandler.EXPECT().Update(gomock.Any(), &patient).Return(&patient, errors.New("id does not exists")),
			},
			expectedError: errors.New("id does not exists"),
		},
	}

	for _, testCase := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := serv.Update(ctx, testCase.input)

		//		_, err := serv.Update(testCase.input)
		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}
