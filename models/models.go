package models

import (
	"fmt"

	"github.com/BimaAdi/fiberGormBoilerplate/settings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)

// Initiate auto migrate the database
func Initiate() {
	fmt.Println("Connect to Database")

	// load from settings
	postgresHost := settings.POSTGRESQL_HOST
	postgresUser := settings.POSTGRESQL_USER
	postgresPassword := settings.POSTGRESQL_PASSWORD
	postgresDatabase := settings.POSTGRESQL_DATABASE
	postgresPort := settings.POSTGRESQL_PORT
	postgresSslMode := settings.POSTGRESQL_SSL_MODE
	postgresTimezone := settings.POSTGRESQL_TIMEZONE

	dsn := "host=" + postgresHost + " user=" + postgresUser + " password=" + postgresPassword +
		" dbname=" + postgresDatabase + " port=" + postgresPort + " sslmode=" + postgresSslMode +
		" TimeZone=" + postgresTimezone

	var err error
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Printf("%#v\n", DBConn)

}

func AutoMigrate() {
	// add models here
	fmt.Println("Migrate Database")
	DBConn.AutoMigrate(&User{})
}

func AutoRollback() {
	fmt.Println("Rollback Database")
	DBConn.Migrator().DropTable(&User{})
}

func ClearAllData() {
	fmt.Println("Clear All Data")
	// DBConn.Exec("DELETE FROM public.oauth2_token")
	// DBConn.Exec("DELETE FROM public.oauth2_session")
	DBConn.Exec("DELETE FROM public.user")
}
