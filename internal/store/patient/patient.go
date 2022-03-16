package store

import (
	"database/sql"
	"strconv"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/anish-kmr/patient-system/internal/model"
)

type PatientStore struct{}

func New() *PatientStore {
	return &PatientStore{}
}

func (e *PatientStore) GetByID(ctx *gofr.Context, id int) (*model.Patient, error) {
	stmt, err := ctx.DB().PrepareContext(
		ctx,
		"SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE id=? AND deletedAt IS NULL",
	)

	if err != nil {
		return nil, errors.Error("internal server error")
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	patient := model.Patient{}

	if err := row.Scan(
		&patient.ID,
		&patient.Name,
		&patient.Phone,
		&patient.Discharged,
		&patient.BloodGroup,
		&patient.Description,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.EntityNotFound{Entity: "Patient", ID: strconv.Itoa(id)}
		}

		return nil, errors.Error("internal server error")
	}

	return &patient, nil
}

func (e *PatientStore) GetAll(ctx *gofr.Context, filters map[string]string) ([]*model.Patient, error) {
	if _, ok := filters["page"]; !ok {
		filters["page"] = "0"
	}

	if _, ok := filters["limit"]; !ok {
		filters["limit"] = "5"
	}

	qry, values := generateFilteredQuery(filters)

	stmt, err := ctx.DB().PrepareContext(ctx, qry)
	if err != nil {
		return nil, errors.Error("internal server error")
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, values...)

	if err != nil {
		return nil, errors.Error("internal server error")
	}
	defer rows.Close()

	patients := []*model.Patient{}

	for rows.Next() {
		patient := model.Patient{}
		if err := rows.Scan(
			&patient.ID,
			&patient.Name,
			&patient.Phone,
			&patient.Discharged,
			&patient.BloodGroup,
			&patient.Description,
			&patient.CreatedAt,
			&patient.UpdatedAt,
		); err != nil {
			return nil, errors.Error("internal server error")
		}

		patients = append(patients, &patient)
	}

	return patients, nil
}

func (e *PatientStore) Create(ctx *gofr.Context, patient *model.Patient) (*model.Patient, error) {
	stmt, err := ctx.DB().PrepareContext(ctx, "INSERT INTO Patient(name,phone,bloodgroup,description) Values(?,?,?,?)")
	if err != nil {
		return nil, errors.Error("internal server error")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(
		ctx,
		patient.Name,
		patient.Phone,
		patient.BloodGroup,
		patient.Description,
	)
	if err != nil {
		return nil, errors.Error("internal server error")
	}

	id, _ := res.LastInsertId()
	newPatient, _ := e.GetByID(ctx, int(id))

	return newPatient, nil
}

func (e *PatientStore) Update(ctx *gofr.Context, id int, patient *model.Patient) (*model.Patient, error) {
	qry, values := generateQuery(patient)
	if len(values) == 0 {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	values = append(values, id)

	stmt, err := ctx.DB().PrepareContext(ctx, qry)

	if err != nil {
		return nil, errors.Error("internal server error")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		return nil, errors.Error("internal server error")
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		return nil, errors.Error("internal server error")
	}

	updatedPatient, _ := e.GetByID(ctx, id)

	return updatedPatient, nil
}

func (e *PatientStore) Delete(ctx *gofr.Context, id int) error {
	stmt, err := ctx.DB().PrepareContext(ctx, "UPDATE Patient SET deletedAt=? WHERE id=? AND deletedAt is NULL")

	if err != nil {
		return errors.Error("internal server error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		time.Now(),
		id,
	)

	if err != nil {
		return errors.Error("internal server error")
	}

	return nil
}
