package stores

import (
	"database/sql"
	"fmt"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/punitj12/patient-app-gofr/internal/models"
)

type DBStorer struct{}

func New() *DBStorer {
	return &DBStorer{}
}

func (dbstore *DBStorer) Create(ctx *gofr.Context, patient *models.Patient) (*models.Patient, error) {
	stmp, err := ctx.DB().
		PrepareContext(ctx, "INSERT INTO `patient`(name,phone,discharged,bloodGroup,description) VALUES(?,?,?,?,?);")
	if err == nil {
		defer stmp.Close()
		res, er := stmp.ExecContext(ctx, patient.Name, patient.Phone, false, patient.BloodGroup, patient.Description)

		if er == nil {
			lastInsertedID, _ := res.LastInsertId()
			patient, err := dbstore.Get(ctx, int(lastInsertedID))

			if err == nil {
				return patient, nil
			}
		}

		return nil, errors.Error("error executing add query")
	}

	return nil, errors.Error("error preparing add query")
}

func (dbstore *DBStorer) Delete(ctx *gofr.Context, id int) error {
	now := time.Now()
	stmp, err := ctx.DB().
		PrepareContext(ctx, "UPDATE `patient` SET deletedAt = ? where id = ? AND deletedAt IS NULL")

	if err == nil {
		defer stmp.Close()
		res, er := stmp.ExecContext(ctx, now, id)

		if er != nil {
			return errors.Error("error executing delete query")
		}

		rowsAff, _ := res.RowsAffected()
		if rowsAff == 1 {
			return nil
		}

		return errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(id)}
	}

	return errors.Error("error preparing delete query")
}

func (dbstore *DBStorer) Get(ctx *gofr.Context, id int) (*models.Patient, error) {
	stmp, err := ctx.DB().
		PrepareContext(ctx, "SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from"+
			" `patient` where id = ? AND deletedAt IS NULL;")

	var patient models.Patient

	if err == nil {
		defer stmp.Close()
		row := stmp.QueryRowContext(ctx, id)
		er := row.
			Scan(&patient.ID, &patient.Name, &patient.Phone, &patient.Discharged,
				&patient.BloodGroup, &patient.Description, &patient.CreatedAt, &patient.UpdatedAt)

		if er != nil {
			if er == sql.ErrNoRows {
				return nil, errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(id)}
			}
		}

		return &patient, nil
	}

	return nil, errors.Error("error preparing select query")
}

func scanRows(rows *sql.Rows) (patients []*models.Patient, err error) {
	for rows.Next() {
		var patient models.Patient
		e := rows.
			Scan(&patient.ID, &patient.Name, &patient.Phone, &patient.Discharged,
				&patient.BloodGroup, &patient.Description, &patient.CreatedAt, &patient.UpdatedAt)

		if e != nil {
			return nil, errors.Error("error scanning data")
		}

		patients = append(patients, &patient)
	}

	return patients, nil
}

func (dbstore *DBStorer) GetAll(ctx *gofr.Context) ([]*models.Patient, error) {
	stmp, err := ctx.DB().
		PrepareContext(ctx, "SELECT id,name,phone,discharged,bloodGroup,description,createdAt,updatedAt from `patient` where deletedAt IS NULL;")

	if err == nil {
		defer stmp.Close()
		rows, er := stmp.QueryContext(ctx)

		if er == nil {
			defer rows.Close()

			patients, err := scanRows(rows)
			if err != nil {
				return nil, errors.Error("error scanning data")
			}

			if len(patients) == 0 {
				return nil, errors.EntityNotFound{Entity: "patient"}
			}

			return patients, nil
		}

		return nil, errors.Error("error fetching data")
	}

	return nil, errors.Error("error preparing select query")
}

func (dbstore *DBStorer) Update(ctx *gofr.Context, patient *models.Patient) (*models.Patient, error) {
	query, values := generateQuery(patient)
	updateQuery := "UPDATE `patient` SET "
	updateQuery = updateQuery + query + " AND deletedAt IS NULL;"
	stmp, err := ctx.DB().PrepareContext(ctx, updateQuery)

	if err == nil {
		defer stmp.Close()
		res, er := stmp.ExecContext(ctx, values...)

		if er == nil {
			defer stmp.Close()

			if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
				newPatient, err := dbstore.Get(ctx, patient.ID)

				if err == nil {
					return newPatient, nil
				}
			}

			return nil, errors.EntityNotFound{Entity: "patient", ID: fmt.Sprint(patient.ID)}
		}

		return nil, errors.Error("error executing update query")
	}

	return nil, errors.Error("error preparing update query")
}
