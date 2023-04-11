package routes

import (
	"github.com/gofiber/fiber/v2"
)

// InitiateRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func InitiateRoutes(app *fiber.App) *fiber.App {
	// noPrefixRoutes := router.Group("/")
	// authRoutes(noPrefixRoutes)
	// userRoutes(noPrefixRoutes)
	authRoutes := app.Group("/auth")
	authRoutes.Post("/login", authLoginRoute)
	authRoutes.Post("/logout", authLogoutRoute)

	// userRoutes := app.Group("/user")
	return app
}
