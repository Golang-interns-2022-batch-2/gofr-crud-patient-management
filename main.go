package main

import (
	"fmt"

	"developer.zopsmart.com/go/gofr/pkg/gofr"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	handlers "github.com/punitj12/patient-app-gofr/internal/http"
	"github.com/punitj12/patient-app-gofr/internal/services"
	"github.com/punitj12/patient-app-gofr/internal/stores"
)

func main() {
	er := godotenv.Load("config/.env")
	if er != nil {
		fmt.Println(er)
		return
	}

	app := gofr.New()
	app.Server.ValidateHeaders = false
	dbstore := stores.New()
	service := services.New(dbstore)
	handle := handlers.New(service)

	app.GET("/patients/{id}", handle.Get)
	app.POST("/patients", handle.Create)
	app.GET("/patients", handle.GetAll)
	app.PUT("/patients/{id}", handle.Update)
	app.DELETE("/patients/{id}", handle.Delete)

	app.Start()
}
