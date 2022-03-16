package store

import (
	"fmt"
	"strconv"

	"github.com/anish-kmr/patient-system/internal/model"
)

func generateQuery(patient *model.Patient) (query string, values []interface{}) {
	query = "UPDATE Patient Set "
	values = []interface{}{}

	if patient.Name.Valid {
		query += "name=? ,"

		values = append(values, patient.Name.String)
	}

	if patient.Phone.Valid {
		query += "phone=? ,"

		values = append(values, patient.Phone.String)
	}

	if patient.BloodGroup.Valid {
		query += "bloodgroup=? ,"

		values = append(values, patient.BloodGroup.String)
	}

	if patient.Discharged.Valid {
		query += "discharged=? ,"

		values = append(values, patient.Discharged.Bool)
	}

	if patient.Description.Valid {
		query += "description=?,"

		values = append(values, patient.Description.String)
	}

	query = query[:len(query)-1]
	query += " WHERE id=? AND deletedAt IS NULL"

	return query, values
}

func generateFilteredQuery(filters map[string]string) (query string, values []interface{}) {
	query = `SELECT id,name,phone,discharged,bloodgroup,description,createdAt,updatedAt FROM Patient WHERE deletedAt IS NULL AND `
	values = []interface{}{}

	for k, v := range filters {
		if k == "page" || k == "limit" {
			continue
		}

		query += fmt.Sprintf("%v=? AND", k)

		values = append(values, v)
	}

	query = query[:len(query)-4]
	query += " LIMIT ?,?"
	page, _ := strconv.Atoi(filters["page"])
	limit, _ := strconv.Atoi(filters["limit"])
	values = append(values, page, limit)

	return query, values
}
