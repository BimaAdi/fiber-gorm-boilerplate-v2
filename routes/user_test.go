package routes_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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

type MigrateTestSuite struct {
	suite.Suite
	app     *fiber.App
	timeout int
}

func (suite *MigrateTestSuite) SetupSuite() {
	settings.InitiateSettings("../.env")
	models.Initiate()
	migrations.MigrateUp("../.env", "file://../migrations/migrations_files/")
	app := fiber.New()
	suite.app = routes.InitiateRoutes(app)
	suite.timeout = 5000 // ms
}

func (suite *MigrateTestSuite) SetupTest() {
	models.ClearAllData()
}

// ==========================================

func (suite *MigrateTestSuite) TestGetPaginateUser() {
	// Given
	timeZoneAsiaJakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err.Error())
	}
	users := []models.User{
		{
			Email:       "a@test.com",
			Username:    "a",
			Password:    "Fakepassword",
			IsActive:    true,
			IsSuperuser: true,
			CreatedAt:   time.Date(2022, 10, 5, 10, 0, 0, 0, timeZoneAsiaJakarta),
		},
		{
			Email:       "b@test.com",
			Username:    "b",
			Password:    "Fakepassword",
			IsActive:    true,
			IsSuperuser: true,
			CreatedAt:   time.Date(2022, 10, 4, 10, 0, 0, 0, timeZoneAsiaJakarta),
		},
		{
			Email:       "c@test.com",
			Username:    "c",
			Password:    "Fakepassword",
			IsActive:    true,
			IsSuperuser: true,
			CreatedAt:   time.Date(2022, 10, 3, 10, 0, 0, 0, timeZoneAsiaJakarta),
		},
	}
	models.DBConn.Create(&users)
	request_user := users[0]
	token, err := core.GenerateJWTTokenFromUser(models.DBConn, request_user)
	if err != nil {
		panic(err.Error())
	}

	// When
	// Test Get Paginate Success
	req, _ := http.NewRequest("GET", "/user/?page=1&page_size=2", nil)
	req.Header.Set("authorization", "Bearer "+token)
	resp, err := suite.app.Test(req, suite.timeout)

	// Expect
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)
	jsonResponse := schemas.UserPaginateResponse{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	err = json.Unmarshal(body, &jsonResponse)
	assert.Nil(suite.T(), err, "Invalid response json")
	assert.Equal(suite.T(), 3, jsonResponse.Counts)
	assert.Equal(suite.T(), 2, jsonResponse.PageCount)
	assert.Equal(suite.T(), 1, jsonResponse.Page)
	assert.Equal(suite.T(), 2, jsonResponse.PageSize)
	assert.Len(suite.T(), jsonResponse.Results, 2)
	for i := 0; i < 2; i++ {
		assert.Equal(suite.T(), users[i].ID, jsonResponse.Results[i].Id)
		assert.Equal(suite.T(), users[i].Username, jsonResponse.Results[i].Username)
		assert.Equal(suite.T(), users[i].Email, jsonResponse.Results[i].Email)
		assert.Equal(suite.T(), users[i].IsActive, jsonResponse.Results[i].IsActive)
	}

	// When 2
	// Test Check Pagination
	req2, _ := http.NewRequest("GET", "/user/?page=2&page_size=2", nil)
	req2.Header.Set("authorization", "Bearer "+token)
	resp2, err := suite.app.Test(req2, suite.timeout)

	// Expect 2
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)
	jsonResponse2 := schemas.UserPaginateResponse{}
	body, err = io.ReadAll(resp2.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	err = json.Unmarshal(body, &jsonResponse2)
	assert.Nil(suite.T(), err, "Invalid response json")
	assert.Equal(suite.T(), 3, jsonResponse2.Counts)
	assert.Equal(suite.T(), 2, jsonResponse2.PageCount)
	assert.Equal(suite.T(), 2, jsonResponse2.Page)
	assert.Equal(suite.T(), 2, jsonResponse2.PageSize)
	assert.Len(suite.T(), jsonResponse2.Results, 1)
	expect_2 := users[2:3]
	for i := 0; i < 1; i++ {
		assert.Equal(suite.T(), expect_2[i].ID, jsonResponse2.Results[i].Id)
		assert.Equal(suite.T(), expect_2[i].Username, jsonResponse2.Results[i].Username)
		assert.Equal(suite.T(), expect_2[i].Email, jsonResponse2.Results[i].Email)
		assert.Equal(suite.T(), expect_2[i].IsActive, jsonResponse2.Results[i].IsActive)
	}

	// When 3
	// Test No Authorization
	req3, _ := http.NewRequest("GET", "/user/?page=1&page_size=2", nil)
	resp3, err := suite.app.Test(req3, suite.timeout)

	// Expect 3
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 401, resp3.StatusCode)
}

func (suite *MigrateTestSuite) TestGetDetailUser() {
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
	// Create user for test
	requestJson := schemas.UserCreateRequest{
		Username:    "test",
		Password:    "testpassword",
		Email:       "test@example.com",
		IsActive:    true,
		IsSuperuser: true,
	}
	requestJsonByte, _ := json.Marshal(requestJson)
	req, _ := http.NewRequest("POST", "/user/", bytes.NewBuffer(requestJsonByte))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "Bearer "+token)
	resp, err := suite.app.Test(req, suite.timeout)
	if err != nil {
		suite.T().Error(err.Error())
	}
	assert.Equal(suite.T(), 201, resp.StatusCode)
	givenJsonResponse := schemas.UserCreateResponse{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	err = json.Unmarshal(body, &givenJsonResponse)
	assert.Nil(suite.T(), err, "Invalid response json")

	// When 1
	// Test success create user
	req1, _ := http.NewRequest("GET", "/user/"+givenJsonResponse.Id, nil)
	req1.Header.Set("authorization", "Bearer "+token)
	resp, err = suite.app.Test(req1, suite.timeout)

	// Expect 1
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)
	jsonResponse1 := schemas.UserDetailResponse{}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	err = json.Unmarshal(body, &jsonResponse1)
	assert.Nil(suite.T(), err, "Invalid response json")

	// When 2
	// Test user not found
	req2, _ := http.NewRequest("GET", "/user/aaaa-bbbbb-ccccc-ddddd", nil)
	req2.Header.Set("authorization", "Bearer "+token)
	resp2, err := suite.app.Test(req2, suite.timeout)

	// Expect 2
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 404, resp2.StatusCode)
	jsonResponse2 := schemas.UserDetailResponse{}
	body, err = io.ReadAll(resp2.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	err = json.Unmarshal(body, &jsonResponse2)
	assert.Nil(suite.T(), err, "Invalid response json")

	// When 3
	// Test No Authorization
	req3, _ := http.NewRequest("GET", "/user/"+givenJsonResponse.Id, nil)
	resp3, err := suite.app.Test(req3, suite.timeout)

	// Expect 3
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 401, resp3.StatusCode)
}

func (suite *MigrateTestSuite) TestCreateUser() {
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
	// Test Create User Success
	requestJson := schemas.UserCreateRequest{
		Username:    "test",
		Password:    "testpassword",
		Email:       "test@example.com",
		IsActive:    true,
		IsSuperuser: true,
	}
	requestJsonByte, _ := json.Marshal(requestJson)
	req, _ := http.NewRequest("POST", "/user/", bytes.NewBuffer(requestJsonByte))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "Bearer "+token)
	resp, err := suite.app.Test(req, suite.timeout)

	// Expect
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode)
	jsonResponse := schemas.UserCreateResponse{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	err = json.Unmarshal(body, &jsonResponse)
	assert.Nil(suite.T(), err, "Invalid response json")

	createdUser := models.User{}
	err = models.DBConn.Where("id = ?", jsonResponse.Id).First(&createdUser).Error
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), createdUser.ID)
	assert.Equal(suite.T(), requestJson.Email, createdUser.Email)
	assert.Equal(suite.T(), requestJson.Username, createdUser.Username)
	assert.Equal(suite.T(), requestJson.IsActive, createdUser.IsActive)
	assert.Equal(suite.T(), requestJson.IsSuperuser, createdUser.IsSuperuser)
	assert.NotNil(suite.T(), createdUser.CreatedAt)

	// When 2
	// Test No Authorization
	requestJson2 := schemas.UserCreateRequest{
		Username:    "test",
		Password:    "testpassword",
		Email:       "test@example.com",
		IsActive:    true,
		IsSuperuser: true,
	}
	requestJsonByte2, _ := json.Marshal(requestJson2)
	req2, _ := http.NewRequest("POST", "/user/", bytes.NewBuffer(requestJsonByte2))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := suite.app.Test(req2, suite.timeout)

	// Expect 2
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 401, resp2.StatusCode)
}

// TODO still error on update user
func (suite *MigrateTestSuite) TestUpdateUser() {
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
	// create user for test
	requestJson := schemas.UserCreateRequest{
		Username:    "test",
		Password:    "testpassword",
		Email:       "test@example.com",
		IsActive:    true,
		IsSuperuser: true,
	}
	requestJsonByte, _ := json.Marshal(requestJson)
	req, _ := http.NewRequest("POST", "/user/", bytes.NewBuffer(requestJsonByte))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "Bearer "+token)
	resp, err := suite.app.Test(req, suite.timeout)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode)
	givenJsonResponse := schemas.UserCreateResponse{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	err = json.Unmarshal(body, &givenJsonResponse)
	assert.Nil(suite.T(), err, "Invalid response json")

	// When 1
	// Test Update User success
	password := "testpassword"
	requestJson1 := schemas.UserUpdateRequest{
		Username:    "test",
		Password:    &password,
		Email:       "test@example.com",
		IsActive:    true,
		IsSuperuser: true,
	}
	requestJsonByte1, _ := json.Marshal(requestJson1)
	req1, _ := http.NewRequest("PUT", "/user/"+givenJsonResponse.Id, bytes.NewBuffer(requestJsonByte1))
	req1.Header.Set("authorization", "Bearer "+token)
	resp1, err := suite.app.Test(req1, suite.timeout)

	// Expect 1
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, resp1.StatusCode)
	jsonResponse1 := schemas.UserUpdateResponse{}
	body, err = io.ReadAll(resp1.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	suite.T().Log(string(body))
	err = json.Unmarshal(body, &jsonResponse1)
	assert.Nil(suite.T(), err, "Invalid response json")

	// When 2
	// Test Update User not found
	password2 := "testpassword2"
	requestJson2 := schemas.UserUpdateRequest{
		Username:    "test",
		Password:    &password2,
		Email:       "test@example.com",
		IsActive:    true,
		IsSuperuser: true,
	}
	requestJsonByte2, _ := json.Marshal(requestJson2)
	req2, _ := http.NewRequest("PUT", "/user/aaaaa-bbbbb-ccccc", bytes.NewBuffer(requestJsonByte2))
	req2.Header.Set("authorization", "Bearer "+token)
	resp2, err := suite.app.Test(req2, suite.timeout)

	// Expect 2
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 404, resp2.StatusCode)
	jsonResponse2 := schemas.NotFoundResponse{}
	body, err = io.ReadAll(resp2.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	err = json.Unmarshal(body, &jsonResponse2)
	assert.Nil(suite.T(), err, "Invalid response json")

	// When 3
	// Test No Authorization
	password3 := "testpassword"
	requestJson3 := schemas.UserUpdateRequest{
		Username:    "test",
		Password:    &password3,
		Email:       "test@example.com",
		IsActive:    true,
		IsSuperuser: true,
	}
	requestJsonByte3, _ := json.Marshal(requestJson3)
	req3, _ := http.NewRequest("PUT", "/user/"+givenJsonResponse.Id, bytes.NewBuffer(requestJsonByte3))
	resp3, err := suite.app.Test(req3, suite.timeout)

	// Expect 3
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 401, resp3.StatusCode)
}

func (suite *MigrateTestSuite) TestDeleteUser() {
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
	// create user for test
	requestJson := schemas.UserCreateRequest{
		Username:    "test",
		Password:    "testpassword",
		Email:       "test@example.com",
		IsActive:    true,
		IsSuperuser: true,
	}
	requestJsonByte, _ := json.Marshal(requestJson)
	req, _ := http.NewRequest("POST", "/user/", bytes.NewBuffer(requestJsonByte))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "Bearer "+token)
	resp, err := suite.app.Test(req, suite.timeout)
	if err != nil {
		suite.T().Error(err.Error())
	}
	assert.Equal(suite.T(), 201, resp.StatusCode)
	jsonResponse := schemas.UserCreateResponse{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Error(err.Error())
	}
	err = json.Unmarshal(body, &jsonResponse)
	assert.Nil(suite.T(), err, "Invalid response json")

	// When 1
	// Test Delete User Success
	req1, _ := http.NewRequest("DELETE", "/user/"+jsonResponse.Id, nil)
	req1.Header.Set("authorization", "Bearer "+token)
	resp1, err := suite.app.Test(req1, suite.timeout)

	// Expect 1
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 204, resp1.StatusCode)

	// When 2
	// Test Delete User Not Found due soft delete
	req2, _ := http.NewRequest("DELETE", "/user/"+jsonResponse.Id, nil)
	req2.Header.Set("authorization", "Bearer "+token)
	resp2, err := suite.app.Test(req2, suite.timeout)

	// Expect 2
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 404, resp2.StatusCode)

	// When 3
	// Test Delete User Not Found
	req3, _ := http.NewRequest("GET", "/user/aaaa-bbbb-cccc", nil)
	req3.Header.Set("authorization", "Bearer "+token)
	resp3, err := suite.app.Test(req3, suite.timeout)

	// Expect 3
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 404, resp3.StatusCode)

	// When 4
	// Test No Authorization
	req4, _ := http.NewRequest("DELETE", "/user/"+jsonResponse.Id, nil)
	resp4, err := suite.app.Test(req4, suite.timeout)

	// Expect 4
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 401, resp4.StatusCode)
}

// ==========================================

func (suite *MigrateTestSuite) TearDownTest() {
	models.ClearAllData()
}

func TestMigrateTestSuite(t *testing.T) {
	suite.Run(t, new(MigrateTestSuite))
}
