package patient

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/shivanisharma200/patient-management/internal/models"
)

type Patient struct {
}

func New() Patient {
	return Patient{}
}
func updateFunc(patient *models.Patient) (query string, values []interface{}) {
	if patient.Name != "" {
		query += "name=?,"

		values = append(values, patient.Name)
	}

	if patient.Description != "" {
		query += "description=?,"

		values = append(values, patient.Description)
	}

	if len(query) > 0 {
		query = query[:len(query)-1]
	}

	return query, values
}

// GetById
func (p Patient) GetByID(ctx *gofr.Context, id int) (*models.Patient, error) {
	const q = "SELECT id,name,phone,discharged,created_at,updated_at,blood_group,description from patients WHERE id = ? AND deleted_at IS NULL"

	patient := models.Patient{}
	err := ctx.DB().QueryRowContext(ctx, q, id).
		Scan(&patient.ID, &patient.Name, &patient.Phone, &patient.Discharged, &patient.CreatedAt, &patient.UpdatedAt,
			&patient.BloodGroup, &patient.Description)

		// sqlerrnorows -- entity not found
		// if err != nil errors.Response 500, message
	if err == sql.ErrNoRows {
		idString := strconv.Itoa(id)
		return nil, errors.EntityNotFound{Entity: "Patient", ID: idString}
	}

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "cannot fetch row",
		}
	}

	return &patient, nil
}

// InsertRow
func (p Patient) Create(ctx *gofr.Context, patient *models.Patient) (*models.Patient, error) {
	resp, err := ctx.DB().ExecContext(ctx, "INSERT INTO patients(name, phone, discharged, blood_group, description) VALUES(?,?,?,?,?)",
		patient.Name, patient.Phone, patient.Discharged, patient.BloodGroup,
		patient.Description)

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "cannot create new patient",
		}
	}

	lastInserted, _ := resp.LastInsertId()

	return p.GetByID(ctx, int(lastInserted))
}

// GetAll
func (p Patient) Get(ctx *gofr.Context) ([]*models.Patient, error) {
	q := "SELECT id, name, phone, discharged, created_at, updated_at, blood_group, description from patients where deleted_at IS NULL"
	rows, err := ctx.DB().
		QueryContext(ctx, q)

	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{Entity: "Patient"}
	}

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "cannot fetch rows",
		}
	}

	var patients []*models.Patient

	defer rows.Close()

	for rows.Next() {
		var patient models.Patient
		_ = rows.Scan(&patient.ID, &patient.Name, &patient.Phone, &patient.Discharged,
			&patient.CreatedAt, &patient.UpdatedAt, &patient.BloodGroup, &patient.Description)

		patients = append(patients, &patient)
	}

	return patients, nil
}

// UpdateById
func (p Patient) Update(ctx *gofr.Context, id int, patient *models.Patient) (*models.Patient, error) {
	query := "UPDATE patients SET "

	resQuery, values := updateFunc(patient)
	query += resQuery
	query += " WHERE id=? AND deleted_at IS NULL"

	values = append(values, id)

	_, err := ctx.DB().ExecContext(ctx, query, values...)

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "cannot update rows",
		}
	}

	return p.GetByID(ctx, id)
}

// DeleteById
func (p Patient) Delete(ctx *gofr.Context, id int) (err error) {
	format := "2006-01-02 15:04:05"

	_, err = ctx.DB().ExecContext(ctx, "UPDATE patients SET deleted_at=? WHERE id=? AND deleted_at IS NULL", time.Now().Format(format), id)

	if err != nil {
		return &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "cannot delete row",
		}
	}

	return nil
}
