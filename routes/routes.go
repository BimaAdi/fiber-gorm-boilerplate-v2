package routes

import (
	"github.com/gofiber/fiber/v2"
)

// InitiateRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func InitiateRoutes(app *fiber.App) *fiber.App {
	authRoutes := app.Group("/auth")
	authRoutes.Post("/login", authLoginRoute)
	authRoutes.Post("/logout", authLogoutRoute)

	userRoutes := app.Group("/user")
	userRoutes.Get("/", GetAllUserRoute)
	userRoutes.Get("/:userId", GetDetailUserRoute)
	userRoutes.Post("/", CreateUserRoute)
	userRoutes.Put("/:userId", UpdateUserRoute)
	userRoutes.Delete("/:userId", DeleteUserRoute)

	return app
}
