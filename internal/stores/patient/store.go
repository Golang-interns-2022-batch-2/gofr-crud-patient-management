package patient

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"errors"
	"fmt"
	"github.com/aakanksha/updated-patient-management-system/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type store struct {
}

func New() *store {
	return &store{}
}

func (s store) GetByID(c *gofr.Context, gid int) (*models.Patient, error) {
	var pt models.Patient
	fmt.Println("start", gid)
	query := "select id,name,phone,discharge,createdat,udatedat,bloodgroup,description from patient where deletedat IS NULL and id=?"

	//fmt.Println("DB ", c.DB().DB.Ping().Error())
	row := c.DB().QueryRowContext(c, query, gid)

	fmt.Println("DEBUG", row)
	err := row.Scan(&pt.Id, &pt.Name, &pt.Phone, &pt.Discharge, &pt.CreatedAt, &pt.UpdatedAt, &pt.BloodGroup, &pt.Description)

	if err != nil {

		return &models.Patient{}, errors.New("id not found")
	}
	fmt.Println(pt)
	return &pt, nil
}

func (s store) Insert(c *gofr.Context, pt *models.Patient) (*models.Patient, error) {
	inserted, err := c.DB().Exec("insert into patient (name,phone,discharge,bloodgroup,description) values (?, ?, ?, ?, ?)",
		pt.Name, pt.Phone, pt.Discharge, pt.BloodGroup, pt.Description)
	if err != nil {
		return &models.Patient{}, errors.New("error in executing insert")
	}
	pid, err := inserted.LastInsertId()
	if err != nil {
		return &models.Patient{}, errors.New("could not get last inserted id")
	} else {
		pt.Id = int(pid)
	}
	return pt, nil
}

func (s *store) Update(c *gofr.Context, pt *models.Patient) (*models.Patient, error) {

	query := "update patient SET name = ?, phone=?, discharge=?,bloodgroup=?,description=? where deletedat IS NULL and id=?"
	_, err := c.DB().Exec(query, &pt.Name, &pt.Phone, &pt.Discharge, &pt.BloodGroup, &pt.Description, &pt.Id)

	if err != nil {
		return &models.Patient{}, errors.New("update failed")
	}

	return pt, nil
}

func (s *store) Delete(c *gofr.Context, did int) error {
	format := "2006-01-02 15:04:05"
	query := "UPDATE patient SET deletedat=? WHERE id=? AND deletedat IS NULL"
	uDeletedAt := time.Now().Format(format)
	res, err := c.DB().Exec(query, uDeletedAt, did)
	if err != nil {
		return err
	}
	rowAffected, _ := res.RowsAffected()
	if rowAffected == 0 {

		return errors.New("No data found")
	}
	return nil
}

func (s *store) GetAll(c *gofr.Context) ([]*models.Patient, error) {
	query := "select id,name,phone,discharge,createdat,udatedat,bloodgroup,description from patient where deletedat IS NULL;"
	var p_array []*models.Patient
	rows, err := c.DB().Query(query)
	defer func() {
		if err == nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		pt := models.Patient{}
		err := rows.Scan(&pt.Id, &pt.Name, &pt.Phone, &pt.Discharge, &pt.CreatedAt, &pt.UpdatedAt, &pt.BloodGroup, &pt.Description)
		if err != nil {
			return nil, err
		}
		p_array = append(p_array, &pt)
	}
	return p_array, nil

}
