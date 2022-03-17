// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package stores is a generated GoMock package.
package stores

import (
	reflect "reflect"

	gofr "developer.zopsmart.com/go/gofr/pkg/gofr"
	models "github.com/aakanksha/updated-patient-management-system/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockStoreInterface is a mock of StoreInterface interface.
type MockStoreInterface struct {
	ctrl     *gomock.Controller
	recorder *MockStoreInterfaceMockRecorder
}

// MockStoreInterfaceMockRecorder is the mock recorder for MockStoreInterface.
type MockStoreInterfaceMockRecorder struct {
	mock *MockStoreInterface
}

// NewMockStoreInterface creates a new mock instance.
func NewMockStoreInterface(ctrl *gomock.Controller) *MockStoreInterface {
	mock := &MockStoreInterface{ctrl: ctrl}
	mock.recorder = &MockStoreInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoreInterface) EXPECT() *MockStoreInterfaceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockStoreInterface) Delete(arg0 *gofr.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStoreInterfaceMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStoreInterface)(nil).Delete), arg0, arg1)
}

// GetAll mocks base method.
func (m *MockStoreInterface) GetAll(arg0 *gofr.Context) ([]*models.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]*models.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockStoreInterfaceMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockStoreInterface)(nil).GetAll), arg0)
}

// GetByID mocks base method.
func (m *MockStoreInterface) GetByID(arg0 *gofr.Context, arg1 int) (*models.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*models.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockStoreInterfaceMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockStoreInterface)(nil).GetByID), arg0, arg1)
}

// Insert mocks base method.
func (m *MockStoreInterface) Insert(arg0 *gofr.Context, arg1 *models.Patient) (*models.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(*models.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockStoreInterfaceMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockStoreInterface)(nil).Insert), arg0, arg1)
}

// Update mocks base method.
func (m *MockStoreInterface) Update(arg0 *gofr.Context, arg1 *models.Patient) (*models.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*models.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockStoreInterfaceMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockStoreInterface)(nil).Update), arg0, arg1)
}
