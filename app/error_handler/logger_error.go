package errorhandler

import (
	"golang-depedency-injection/app/model/response"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {

	// message := err.Error()
	// logStop := util.LogResponse(ctx, message, "")
	// fmt.Println(logStop)

	_, databaseError := err.(DatabaseError)
	if databaseError {
		return ctx.Status(500).JSON(response.GeneralResponse{
			StatusCode: 500,
			Message:    "Terjadi kesalahan, silahkan ulangi beberapa saat lagi, code=[001]",
		})
	}

	_, dataNotFoundError := err.(DataNotFoundError)
	if dataNotFoundError {
		return ctx.Status(404).JSON(response.GeneralResponse{
			StatusCode: 404,
			Message:    err.Error(),
		})
	}

	_, generalError := err.(GeneralError)
	if generalError {
		return ctx.Status(400).JSON(response.GeneralResponse{
			StatusCode: 400,
			Message:    err.Error(),
		})
	}

	return ctx.Status(400).JSON(response.GeneralResponse{
		StatusCode: 400,
		// Message:    "Terjadi kesalahan, silahkan ulangi beberapa saat lagi, code=[002]",
		Message: err.Error(),
	})
}
