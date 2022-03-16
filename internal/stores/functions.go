package stores

import "github.com/punitj12/patient-app-gofr/internal/models"

func generateQuery(patient *models.Patient) (query string, values []interface{}) {
	if patient.Name.Valid {
		query += " name = ?, "

		values = append(values, patient.Name.String)
	}

	if patient.BloodGroup.Valid {
		query += " bloodGroup = ?, "

		values = append(values, patient.BloodGroup.String)
	}

	if patient.Description.Valid {
		query += " description = ?, "

		values = append(values, patient.Description.String)
	}

	if patient.Discharged.Valid {
		query += " discharged = ?, "

		values = append(values, patient.Discharged.Bool)
	}

	if patient.Phone.Valid {
		query += " phone = ?, "

		values = append(values, patient.Phone.String)
	}

	query = query[:len(query)-2] + " "
	query += "where id = ?"

	values = append(values, patient.ID)

	return query, values
}
