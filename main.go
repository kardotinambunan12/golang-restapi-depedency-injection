package main

import (
	"golang-depedency-injection/app/config"
	"golang-depedency-injection/app/controller"
	errorhandler "golang-depedency-injection/app/error_handler"
	"golang-depedency-injection/app/repository"
	"golang-depedency-injection/app/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// configuration := config.New()

	// setup repository
	testApiRepository := repository.NewTestApiRepository()

	// setup service
	testApiService := service.NewTestApiService(&testApiRepository)

	// setup controller
	testApiController := controller.NewTestRestApiController(&testApiService)

	app := fiber.New(config.NewFiberConfig())
	app.Use(recover.New())

	app.Use(logger.New())

	// setup routing
	testApiController.Route(app)

	err := app.Listen(":3001")
	errorhandler.PanicIfNeeded(err)

}
