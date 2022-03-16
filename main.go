package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	httpPatient "github.com/shivanisharma200/patient-management/internal/http/patient"
	servicePatient "github.com/shivanisharma200/patient-management/internal/service/patient"
	storePatient "github.com/shivanisharma200/patient-management/internal/store/patient"
)

func main() {
	store := storePatient.New()

	service := servicePatient.New(store)

	postHandler := httpPatient.New(service)

	g := gofr.New()

	// REST-API
	// creating edpoints
	g.Server.ValidateHeaders = false
	g.GET("/patients/{id}", postHandler.GetByID)
	g.POST("/patients", postHandler.Create)
	g.GET("/patients", postHandler.Get)
	g.PUT("/patients/{id}", postHandler.Update)
	g.DELETE("/patients/{id}", postHandler.Delete)
	g.Start()
}
