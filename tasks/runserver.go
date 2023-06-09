package tasks

import (
	_ "github.com/BimaAdi/fiberGormBoilerplate/docs"
	"github.com/BimaAdi/fiberGormBoilerplate/models"
	"github.com/BimaAdi/fiberGormBoilerplate/routes"
	"github.com/BimaAdi/fiberGormBoilerplate/settings"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func RunServer(envPath string) {
	// Initialize environtment variable
	settings.InitiateSettings(envPath)

	// Initiate Database connection
	models.Initiate()

	// development or release
	// if settings.GIN_MODE == "release" {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	// Initiate fiber app
	app := fiber.New()

	// Cors Middleware
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))

	// TODO Initiate static and template
	// router.Static("/assets", "./assets")
	// router.LoadHTMLGlob("templates/*.html")

	// Initialize fiber route
	app = routes.InitiateRoutes(app)

	// setup swagger
	app.Get("/docs/*", swagger.HandlerDefault)

	// run fiber server
	app.Listen(settings.SERVER_HOST + ":" + settings.SERVER_PORT)
}
