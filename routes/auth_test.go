package routes_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/BimaAdi/fiberGormBoilerplate/core"
	"github.com/BimaAdi/fiberGormBoilerplate/migrations"
	"github.com/BimaAdi/fiberGormBoilerplate/models"
	"github.com/BimaAdi/fiberGormBoilerplate/routes"
	"github.com/BimaAdi/fiberGormBoilerplate/schemas"
	"github.com/BimaAdi/fiberGormBoilerplate/settings"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MigrateAuthTestSuite struct {
	suite.Suite
	app     *fiber.App
	timeout int
}

func (suite *MigrateAuthTestSuite) SetupSuite() {
	settings.InitiateSettings("../.env")
	models.Initiate()
	migrations.MigrateUp("../.env", "file://../migrations/migrations_files/")
	app := fiber.New()
	suite.app = routes.InitiateRoutes(app)
	suite.timeout = 5 // second
}

func (suite *MigrateAuthTestSuite) SetupTest() {
	models.ClearAllData()
}

func (suite *MigrateAuthTestSuite) TestLoginSuccess() {
	// Given
	timeZoneAsiaJakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err.Error())
	}
	hashPasword, err := core.HashPassword("Fakepassword")
	if err != nil {
		panic(err.Error())
	}
	user_login := models.User{
		Email:       "test@test.com",
		Username:    "test",
		Password:    hashPasword,
		IsActive:    true,
		IsSuperuser: true,
		CreatedAt:   time.Date(2022, 10, 5, 10, 0, 0, 0, timeZoneAsiaJakarta),
	}
	models.DBConn.Create(&user_login)

	// When
	var param = url.Values{}
	param.Set("username", "test")
	param.Set("password", "Fakepassword")
	var payload = bytes.NewBufferString(param.Encode())
	req, _ := http.NewRequest("POST", "/auth/login", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := suite.app.Test(req, suite.timeout)

	// Expect
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)
	jsonResponse := schemas.LoginResponse{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	err = json.Unmarshal(body, &jsonResponse)
	assert.Nil(suite.T(), err, "Invalid response json")
}

func (suite *MigrateAuthTestSuite) TestLoginFailed() {
	// Given
	timeZoneAsiaJakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err.Error())
	}
	hashPasword, err := core.HashPassword("Fakepassword")
	if err != nil {
		panic(err.Error())
	}
	user_login := models.User{
		Email:       "test@test.com",
		Username:    "test",
		Password:    hashPasword,
		IsActive:    true,
		IsSuperuser: true,
		CreatedAt:   time.Date(2022, 10, 5, 10, 0, 0, 0, timeZoneAsiaJakarta),
	}
	models.DBConn.Create(&user_login)

	// When
	var param = url.Values{}
	param.Set("username", "test")
	param.Set("password", "wrong password")
	var payload = bytes.NewBufferString(param.Encode())
	req, _ := http.NewRequest("POST", "/auth/login", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := suite.app.Test(req, suite.timeout)

	// Expect
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 400, resp.StatusCode)
}

func (suite *MigrateAuthTestSuite) TestLogoutSuccess() {
	// Given
	// create request user
	timeZoneAsiaJakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err.Error())
	}
	request_user := models.User{
		Email:       "a@test.com",
		Username:    "a",
		Password:    "Fakepassword",
		IsActive:    true,
		IsSuperuser: true,
		CreatedAt:   time.Date(2022, 10, 5, 10, 0, 0, 0, timeZoneAsiaJakarta),
	}
	models.DBConn.Create(&request_user)
	token, err := core.GenerateJWTTokenFromUser(models.DBConn, request_user)
	if err != nil {
		panic(err.Error())
	}

	// When
	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("authorization", "Bearer "+token)
	resp, err := suite.app.Test(req, suite.timeout)

	// Expect
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	jsonResponse := schemas.LogoutResponse{}
	err = json.Unmarshal(body, &jsonResponse)
	assert.Nil(suite.T(), err, "Invalid response json")
}

func (suite *MigrateAuthTestSuite) TestLogoutInvalidToken() {
	// Given
	token := "theinvalidtoken"

	// When
	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("authorization", "Bearer "+token)
	resp, err := suite.app.Test(req, suite.timeout)

	// Expect
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 401, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	jsonResponse := schemas.UnauthorizedResponse{}
	err = json.Unmarshal(body, &jsonResponse)
	assert.Nil(suite.T(), err, "Invalid response json")
}

func (suite *MigrateAuthTestSuite) TearDownTest() {
	models.ClearAllData()
}

func TestMigrateAuthTestSuite(t *testing.T) {
	suite.Run(t, new(MigrateAuthTestSuite))
}
