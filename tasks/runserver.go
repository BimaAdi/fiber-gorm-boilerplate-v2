package tasks

import (
	"github.com/BimaAdi/fiberGormBoilerplate/models"
	"github.com/BimaAdi/fiberGormBoilerplate/routes"
	"github.com/BimaAdi/fiberGormBoilerplate/settings"
	"github.com/gofiber/fiber/v2"
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

	// TODO Cors Middleware
	// router := gin.Default()
	// router.Use(cors.New(cors.Config{
	// 	AllowAllOrigins:        true,
	// 	AllowOrigins:           []string{},
	// 	AllowMethods:           []string{"GET", "POST", "PUT", "DELETE", "OPTION"},
	// 	AllowHeaders:           []string{"Origin", "Content-Type", "authorization", "accept"},
	// 	AllowCredentials:       true,
	// 	ExposeHeaders:          []string{"Content-Length"},
	// 	MaxAge:                 0,
	// 	AllowWildcard:          true,
	// 	AllowBrowserExtensions: true,
	// 	AllowWebSockets:        true,
	// 	AllowFiles:             true,
	// }))

	// TODO Initiate static and template
	// router.Static("/assets", "./assets")
	// router.LoadHTMLGlob("templates/*.html")

	// Initialize fiber route
	// routes := routes.GetRoutes(router)
	app := fiber.New()
	app = routes.InitiateRoutes(app)

	// TODO setup swagger
	// docs.SwaggerInfo.BasePath = "/"
	// docs.SwaggerInfo.Host = settings.SERVER_HOST + ":" + settings.SERVER_PORT
	// routes.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// TODO run fiber server
	// routes.Run(settings.SERVER_HOST + ":" + settings.SERVER_PORT)
	app.Listen(settings.SERVER_HOST + ":" + settings.SERVER_PORT)
}
