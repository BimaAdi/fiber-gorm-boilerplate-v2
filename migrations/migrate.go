package migrations

import (
	"errors"
	"strconv"

	"github.com/BimaAdi/fiberGormBoilerplate/settings"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Up(file_path string, postgres_url string, step *int) {
	m, err := migrate.New(
		file_path,
		postgres_url)
	if err != nil {
		panic(err.Error())
	}

	if step == nil {
		m.Up()
	} else {
		m.Steps(*step)
	}
}

func Down(file_path string, postgres_url string, step *int) {
	m, err := migrate.New(
		file_path,
		postgres_url)
	if err != nil {
		panic(err.Error())
	}

	if step == nil {
		m.Down()
	} else {
		m.Steps(*step)
	}
}

func ShowCurrentVersion(file_path string, postgres_url string) string {
	m, err := migrate.New(
		file_path,
		postgres_url)
	if err != nil {
		panic(err.Error())
	}

	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return "No migration applied"
		}
		panic(err.Error())
	}
	dirtyString := ""
	if dirty {
		dirtyString = "true"
	} else {
		dirtyString = "false"
	}
	return "version: " + strconv.FormatUint(uint64(version), 10) + " dirty: " + dirtyString
}

func MigrateUp(envPath string, migrations_file_path string) {
	// Initialize environtment variable
	settings.InitiateSettings(envPath)
	postgresUser := settings.POSTGRESQL_USER
	postgresPassword := settings.POSTGRESQL_PASSWORD
	postgresHost := settings.POSTGRESQL_HOST
	postgresPort := settings.POSTGRESQL_PORT
	postgresDatabase := settings.POSTGRESQL_DATABASE
	postgresSslMode := settings.POSTGRESQL_SSL_MODE

	// Initiate Database connection
	postgres_url := "postgres://" + postgresUser + ":" + postgresPassword + "@" +
		postgresHost + ":" + postgresPort + "/" + postgresDatabase + "?sslmode=" + postgresSslMode

	// Migrate
	Up(migrations_file_path, postgres_url, nil)
}

func MigrateStep(envPath string, migrations_file_path string, step *int) {
	// Initialize environtment variable
	settings.InitiateSettings(envPath)
	postgresUser := settings.POSTGRESQL_USER
	postgresPassword := settings.POSTGRESQL_PASSWORD
	postgresHost := settings.POSTGRESQL_HOST
	postgresPort := settings.POSTGRESQL_PORT
	postgresDatabase := settings.POSTGRESQL_DATABASE
	postgresSslMode := settings.POSTGRESQL_SSL_MODE

	// Initiate Database connection
	postgres_url := "postgres://" + postgresUser + ":" + postgresPassword + "@" +
		postgresHost + ":" + postgresPort + "/" + postgresDatabase + "?sslmode=" + postgresSslMode

	// Migrate
	Up(migrations_file_path, postgres_url, step)
}

func MigrateDown(envPath string, migrations_file_path string) {
	// Initialize environtment variable
	settings.InitiateSettings(envPath)
	postgresUser := settings.POSTGRESQL_USER
	postgresPassword := settings.POSTGRESQL_PASSWORD
	postgresHost := settings.POSTGRESQL_HOST
	postgresPort := settings.POSTGRESQL_PORT
	postgresDatabase := settings.POSTGRESQL_DATABASE
	postgresSslMode := settings.POSTGRESQL_SSL_MODE

	// Initiate Database connection
	postgres_url := "postgres://" + postgresUser + ":" + postgresPassword + "@" +
		postgresHost + ":" + postgresPort + "/" + postgresDatabase + "?sslmode=" + postgresSslMode

	// Migrate
	Down(migrations_file_path, postgres_url, nil)
}

func GetMigrateVersion(envPath string, migrations_file_path string) string {
	// Initialize environtment variable
	settings.InitiateSettings(envPath)
	postgresUser := settings.POSTGRESQL_USER
	postgresPassword := settings.POSTGRESQL_PASSWORD
	postgresHost := settings.POSTGRESQL_HOST
	postgresPort := settings.POSTGRESQL_PORT
	postgresDatabase := settings.POSTGRESQL_DATABASE
	postgresSslMode := settings.POSTGRESQL_SSL_MODE

	// Initiate Database connection
	postgres_url := "postgres://" + postgresUser + ":" + postgresPassword + "@" +
		postgresHost + ":" + postgresPort + "/" + postgresDatabase + "?sslmode=" + postgresSslMode

	return ShowCurrentVersion(migrations_file_path, postgres_url)
}
