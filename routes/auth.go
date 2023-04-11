package routes

import (
	"github.com/BimaAdi/fiberGormBoilerplate/core"
	"github.com/BimaAdi/fiberGormBoilerplate/models"
	"github.com/BimaAdi/fiberGormBoilerplate/repository"
	"github.com/BimaAdi/fiberGormBoilerplate/schemas"
	"github.com/gofiber/fiber/v2"
)

// Login
//
//	@Summary		Login
//	@Description	login
//	@Tags			Auth
//	@Produce		json
//	@Param			payload	formData	schemas.LoginFormRequest	true	"form data"
//	@Success		200		{object}	schemas.LoginResponse
//	@Failure		400		{object}	schemas.BadRequestResponse
//	@Failure		500		{object}	schemas.InternalServerErrorResponse
//	@Router			/auth/login [post]
func authLoginRoute(c *fiber.Ctx) error {
	// Get data from form
	formRequest := schemas.LoginFormRequest{}
	if err := c.BodyParser(&formRequest); err != nil {
		return c.Status(400).JSON(schemas.BadRequestResponse{
			Message: err.Error(),
		})
	}

	// Get User
	user, err := repository.GetUserByUsername(models.DBConn, formRequest.Username)
	if err != nil {
		return c.Status(400).JSON(schemas.BadRequestResponse{
			Message: "invalid credentials",
		})
	}

	// Check Password
	if !core.CheckPasswordHash(formRequest.Password, user.Password) {
		return c.Status(400).JSON(schemas.BadRequestResponse{
			Message: "invalid credentials",
		})
	}

	// Generate JWT token
	token, err := core.GenerateJWTTokenFromUser(models.DBConn, user)
	if err != nil {
		return c.Status(500).JSON(schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(200).JSON(schemas.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}

// Logout
//
//	@Summary		Logout
//	@Description	logout
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	schemas.LogoutResponse
//	@Failure		400	{object}	schemas.UnauthorizedResponse
//	@Failure		500	{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/auth/logout [post]
func authLogoutRoute(c *fiber.Ctx) error {
	// Authorize User
	user, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		return c.Status(401).JSON(schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
	}

	return c.Status(200).JSON(schemas.LogoutResponse{
		Email:    user.Email,
		Username: user.Username,
	})
}
