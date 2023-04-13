package routes

import (
	"errors"
	"time"

	"github.com/BimaAdi/fiberGormBoilerplate/core"
	"github.com/BimaAdi/fiberGormBoilerplate/models"
	"github.com/BimaAdi/fiberGormBoilerplate/repository"
	"github.com/BimaAdi/fiberGormBoilerplate/schemas"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Get All User
//
//	@Summary		Get All User
//	@Description	Get All User
//	@Tags			User
//	@Produce		json
//	@Param			page		query		int	false	"page"
//	@Param			page_size	query		int	false	"page"
//	@Success		200			{object}	schemas.UserPaginateResponse
//	@Failure		400			{object}	schemas.BadRequestResponse
//	@Failure		401			{object}	schemas.UnauthorizedResponse
//	@Failure		500			{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/user/ [get]
func GetAllUserRoute(c *fiber.Ctx) error {
	// Authorize User
	_, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		return c.Status(401).JSON(schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
	}

	// Get Query Parameter
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	search := c.Query("search", "")
	if page <= 0 || pageSize <= 0 {
		errorResponse := []map[string]string{}
		if page <= 0 {
			x := map[string]string{
				"page": "invalid page, page should positive integer",
			}
			errorResponse = append(errorResponse, x)
		}

		if pageSize <= 0 {
			x := map[string]string{
				"page_size": "invalid page_size, page_size should positive integer",
			}
			errorResponse = append(errorResponse, x)
		}

		return c.Status(422).JSON(schemas.UnprocessableEntityResponse{
			Message: errorResponse,
		})
	}
	var searchNilable *string = nil
	if search != "" {
		searchNilable = &search
	}

	users, numData, numPage, err := repository.GetPaginatedUser(
		models.DBConn, page, pageSize, searchNilable,
	)
	if err != nil {
		return c.Status(500).JSON(schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
	}

	arrayDetailUser := []schemas.UserDetailResponse{}
	for _, item := range users {
		arrayDetailUser = append(arrayDetailUser, schemas.UserDetailResponse{
			Id:       item.ID,
			Username: item.Username,
			Email:    item.Email,
			IsActive: item.IsActive,
		})
	}

	return c.Status(200).JSON(schemas.UserPaginateResponse{
		Counts:    int(numData),
		PageCount: int(numPage),
		PageSize:  pageSize,
		Page:      page,
		Results:   arrayDetailUser,
	})
}

// Get Detail User
//
//	@Summary		Get Detail User
//	@Description	Get detail user
//	@Tags			User
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	schemas.UserDetailResponse
//	@Failure		400	{object}	schemas.BadRequestResponse
//	@Failure		404	{object}	schemas.NotFoundResponse
//	@Failure		500	{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/user/{id} [get]
func GetDetailUserRoute(c *fiber.Ctx) error {
	// Authorize User
	_, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		return c.Status(401).JSON(schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
	}

	// Get Params
	userId := c.Params("userId")
	if !core.IsValidUUID(userId) {
		return c.Status(404).JSON(schemas.NotFoundResponse{
			Message: "user not found",
		})
	}

	user, err := repository.GetUserById(models.DBConn, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(schemas.NotFoundResponse{
				Message: "user not found",
			})
		}
		return c.Status(500).JSON(schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(200).JSON(schemas.UserDetailResponse{
		Id:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		IsActive:    user.IsActive,
		IsSuperuser: user.IsSuperuser,
	})
}

// Create User
//
//	@Summary		Create User
//	@Description	Create User
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		schemas.UserCreateRequest	true	"Create User"
//	@Success		200		{object}	schemas.UserCreateResponse
//	@Failure		400		{object}	schemas.BadRequestResponse
//	@Failure		422		{object}	schemas.UnprocessableEntityResponse
//	@Failure		500		{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/user/ [post]
func CreateUserRoute(c *fiber.Ctx) error {
	// Authorize User
	_, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		return c.Status(401).JSON(schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
	}

	// validation
	var newUser schemas.UserCreateRequest
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(400).JSON(schemas.BadRequestResponse{
			Message: err.Error(),
		})
	}
	is_valid, validation_errors := core.ValidateSchemas(newUser)
	if !is_valid {
		return c.Status(422).JSON(validation_errors)
	}

	now := time.Now()
	createdUser, err := repository.CreateUser(
		models.DBConn,
		newUser.Username,
		newUser.Email,
		newUser.Password,
		newUser.IsActive,
		newUser.IsSuperuser,
		now,
		&now,
	)
	if err != nil {
		return c.Status(500).JSON(schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(201).JSON(schemas.UserCreateResponse{
		Id:          createdUser.ID,
		Username:    createdUser.Username,
		Email:       createdUser.Email,
		IsActive:    createdUser.IsActive,
		IsSuperuser: createdUser.IsSuperuser,
	})
}

// Update User
//
//	@Summary		Update User
//	@Description	Update User
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"User ID"
//	@Param			user	body		schemas.UserUpdateRequest	true	"Update User"
//	@Success		200		{object}	schemas.UserUpdateResponse
//	@Failure		400		{object}	schemas.BadRequestResponse
//	@Failure		404		{object}	schemas.NotFoundResponse
//	@Failure		422		{object}	schemas.UnprocessableEntityResponse
//	@Failure		500		{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/user/{id} [put]
func UpdateUserRoute(c *fiber.Ctx) error {
	// Authorize User
	_, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		return c.Status(401).JSON(schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
	}

	// get input user
	userId := c.Params("userId")
	if !core.IsValidUUID(userId) {
		return c.Status(404).JSON(schemas.NotFoundResponse{
			Message: "user not found",
		})
	}

	// validation
	jsonRequest := schemas.UserUpdateRequest{}
	if err = c.BodyParser(&jsonRequest); err != nil {
		return c.Status(400).JSON(schemas.BadRequestResponse{
			Message: err.Error(),
		})
	}
	is_valid, validation_errors := core.ValidateSchemas(jsonRequest)
	if !is_valid {
		return c.Status(422).JSON(validation_errors)
	}

	// get existing user
	user, err := repository.GetUserById(models.DBConn, userId)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(schemas.NotFoundResponse{
				Message: "user not found",
			})
		}
		return c.Status(500).JSON(schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
	}

	// update user
	updatedUser, err := repository.UpdateUser(
		models.DBConn,
		user,
		jsonRequest.Email,
		jsonRequest.Username,
		jsonRequest.Password,
		jsonRequest.IsActive,
		jsonRequest.IsSuperuser,
	)

	if err != nil {
		return c.Status(500).JSON(schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(200).JSON(schemas.UserUpdateResponse{
		Id:          updatedUser.ID,
		Username:    updatedUser.Username,
		Email:       updatedUser.Email,
		IsActive:    updatedUser.IsActive,
		IsSuperuser: updatedUser.IsSuperuser,
	})
}

// Delete User
//
//	@Summary		Delete User
//	@Description	Delete user
//	@Tags			User
//	@Param			id	path	string	true	"User ID"
//	@Success		204
//	@Failure		404	{object}	schemas.NotFoundResponse
//	@Failure		500	{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/user/{id} [delete]
func DeleteUserRoute(c *fiber.Ctx) error {
	// Authorize User
	_, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		return c.Status(401).JSON(schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
	}

	// get input user
	userId := c.Params("userId")
	if !core.IsValidUUID(userId) {
		return c.Status(404).JSON(schemas.NotFoundResponse{
			Message: "user not found",
		})
	}

	// get existing user
	user, err := repository.GetUserById(models.DBConn, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(schemas.NotFoundResponse{
				Message: "user not found",
			})
		}
		return c.Status(500).JSON(schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
	}

	_, err = repository.DeleteUser(models.DBConn, user)
	if err != nil {
		return c.Status(500).JSON(schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
	}
	return c.Status(204).JSON(nil)
}
