package main

import (
	http "github.com/aakanksha/updated-patient-management-system/internal/http/patient"
	svc "github.com/aakanksha/updated-patient-management-system/internal/service/patient"
	str "github.com/aakanksha/updated-patient-management-system/internal/stores/patient"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	//"github.com/aakanksha/updated-patient-management-system/internal/configs"
)

func main() {
	app := gofr.New()
	storelevel := str.New()
	servicelevel := svc.New(storelevel)
	handlerlevel := http.New(servicelevel)

	app.Server.ValidateHeaders = false
	app.Server.HTTP.Port = 10000
	//app.Server.HTTP.RedirectToHTTPS = false

	app.GET("/patients/{id}", handlerlevel.GetByID)
	app.POST("/patients", handlerlevel.Insert)
	app.PUT("/patients/{id}", handlerlevel.Update)
	app.GET("/patients", handlerlevel.GetAll)
	app.DELETE("/patients/{id}", handlerlevel.Delete)
	app.Start()
}
