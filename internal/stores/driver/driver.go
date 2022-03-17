package driver

import (
	"database/sql"
	"fmt"
)

func Connection() *sql.DB {
	db, err := sql.Open("mysql", "root:<password>@tcp(127.0.0.1:3306)/patientDB?parseTime=true")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Connected")
	}
	return db
}
