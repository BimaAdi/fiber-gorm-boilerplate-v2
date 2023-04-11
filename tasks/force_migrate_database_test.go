package tasks_test

import (
	"testing"

	"github.com/BimaAdi/fiberGormBoilerplate/migrations"
	"github.com/BimaAdi/fiberGormBoilerplate/tasks"
	"github.com/stretchr/testify/assert"
)

func TestForceMigrate(t *testing.T) {
	// Given
	migrations.MigrateDown("../.env", "file://../migrations/migrations_files/")
	migrations.MigrateUp("../.env", "file://../migrations/migrations_files/")

	// When
	tasks.ForceRollback("../.env")
	tasks.ForceMigrate("../.env")
	tasks.ForceRollback("../.env")
	assert.Equal(
		t, "No migration applied",
		migrations.GetMigrateVersion("../.env", "file://../migrations/migrations_files/"),
	)
	migrations.MigrateUp("../.env", "file://../migrations/migrations_files/")
	assert.NotEqual(
		t, "No migration applied",
		migrations.GetMigrateVersion("../.env", "file://../migrations/migrations_files/"),
	)
	migrations.MigrateDown("../.env", "file://../migrations/migrations_files/")
	assert.Equal(
		t, "No migration applied",
		migrations.GetMigrateVersion("../.env", "file://../migrations/migrations_files/"),
	)

	// Expect
	// No Error
}
