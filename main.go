package main

import (
	"fmt"
	"os"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	handler "github.com/anish-kmr/patient-system/internal/http/patient"
	service "github.com/anish-kmr/patient-system/internal/service/patient"
	store "github.com/anish-kmr/patient-system/internal/store/patient"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		fmt.Println("Error loading .env file")
	}

	fmt.Println(os.Getenv("DB_HOST"))

	patientStore := store.New()
	patientService := service.New(patientStore)
	patientHandler := handler.New(patientService)

	app := gofr.New()
	app.Server.ValidateHeaders = false

	app.GET("/patient/{id}", patientHandler.GetByID)
	app.GET("/patient", patientHandler.GetAll)
	app.POST("/patient", patientHandler.Create)
	app.PUT("/patient/{id}", patientHandler.Update)
	app.DELETE("/patient/{id}", patientHandler.Delete)

	app.Start()
}
