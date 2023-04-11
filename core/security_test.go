package core_test

import (
	"testing"

	"github.com/BimaAdi/fiberGormBoilerplate/core"
	"github.com/BimaAdi/fiberGormBoilerplate/migrations"
	"github.com/BimaAdi/fiberGormBoilerplate/models"
	"github.com/BimaAdi/fiberGormBoilerplate/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestGenerateJWTToken(t *testing.T) {
	settings.InitiateSettings("../.env")
	token, err := core.GenerateJWTToken("aaaaa-bbbbb-ccccc", "bimaadi419@gmail.com")
	assert.Nil(t, err)
	id, email, err := core.GetPayloadFromJWTToken(token)
	assert.Equal(t, "aaaaa-bbbbb-ccccc", id)
	assert.Equal(t, "bimaadi419@gmail.com", email)
}

func TestExpiredJWTToken(t *testing.T) {
	settings.InitiateSettings("../.env")
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImJpbWFhZGk0MTlAZ21haWwuY29tIiwiZXhwIjoxNjc2MjYwMDg0LCJpYXQiOjE2NzYyNjE4ODQsImlkIjoiYWFhYWEtYmJiYmItY2NjY2MifQ.8VwyTsDIfUf8WT1A9MxFHMqWn-ZVi9RZACvA2y5K9WE"
	_, _, err := core.GetPayloadFromJWTToken(token)
	assert.NotNil(t, err)
}

type MigrateTestSuite struct {
	suite.Suite
}

func (suite *MigrateTestSuite) SetupSuite() {
	settings.InitiateSettings("../.env")
	models.Initiate()
	migrations.MigrateUp("../.env", "file://../migrations/migrations_files/")
}

func (suite *MigrateTestSuite) SetupTest() {
	models.ClearAllData()
}

// ==========================================

func (suite *MigrateTestSuite) TestGetUserFromJwtToken() {
	// Given
	user := models.User{
		Email:       "bimaadi419@gmail.com",
		Password:    "hashpassword",
		IsActive:    true,
		IsSuperuser: true,
	}
	models.DBConn.Create(&user)

	// When
	token, err := core.GenerateJWTTokenFromUser(models.DBConn, user)
	assert.Nil(suite.T(), err)
	userPayload, err := core.GetUserFromJWTToken(models.DBConn, token)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), user.ID, userPayload.ID)
	assert.Equal(suite.T(), user.Email, userPayload.Email)
	assert.Equal(suite.T(), user.IsActive, userPayload.IsActive)
	assert.Equal(suite.T(), user.IsSuperuser, userPayload.IsSuperuser)
}

func (suite *MigrateTestSuite) TestGetUserFromJwtTokenInvalidToken() {
	// Given
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImJpbWFhZGk0MTlAZ21haWwuY29tIiwiZXhwIjoxNjc2MjYwMDg0LCJpYXQiOjE2NzYyNjE4ODQsImlkIjoiYWFhYWEtYmJiYmItY2NjY2MifQ.8VwyTsDIfUf8WT1A9MxFHMqWn-ZVi9RZACvA2y5K9WE"

	// When
	_, err := core.GetUserFromJWTToken(models.DBConn, token)

	// Expect
	assert.NotNil(suite.T(), err)
}

// ==========================================

func (suite *MigrateTestSuite) TearDownTest() {
	models.ClearAllData()
}

func TestMigrateTestSuite(t *testing.T) {
	suite.Run(t, new(MigrateTestSuite))
}
