package tasks

import (
	"time"

	"github.com/BimaAdi/fiberGormBoilerplate/models"
	"github.com/BimaAdi/fiberGormBoilerplate/repository"
	"github.com/BimaAdi/fiberGormBoilerplate/settings"
)

func CreateSuperUser(envPath string, email string, username string, password string) {
	// Initialize environtment variable
	settings.InitiateSettings(envPath)

	// Initiate Database connection
	models.Initiate()

	now := time.Now()
	repository.CreateUser(models.DBConn, username, email, password, true, true, now, &now)
}
