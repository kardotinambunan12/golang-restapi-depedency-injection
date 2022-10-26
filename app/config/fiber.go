package config

import (
	errorhandler "golang-depedency-injection/app/error_handler"

	"github.com/gofiber/fiber/v2"
)

func NewFiberConfig() fiber.Config {

	return fiber.Config{
		ErrorHandler: errorhandler.ErrorHandler,
	}
}
