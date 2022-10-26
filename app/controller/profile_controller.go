package controller

import (
	"encoding/json"
	errorhandler "golang-depedency-injection/app/error_handler"
	"golang-depedency-injection/app/model/request"
	"golang-depedency-injection/app/service"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TestRestApiController struct {
	TestRestApiService service.TestRestApiService
}

func NewTestRestApiController(testRestApiService *service.TestRestApiService) TestRestApiController {
	return TestRestApiController{TestRestApiService: *testRestApiService}
}

func (controller *TestRestApiController) Route(app *fiber.App) {
	app.Post("/api/profile/GetAll", controller.TestApiGetList)
	app.Get("/api/profile/GetbyId/:profileId", controller.TestAPiByID)
	app.Post("/api/profile/Insert", controller.InsertProfile)
	app.Put("/api/profile/Updated/:profileId", controller.Updatedprofile)
	app.Delete("/api/profile/Deleted/:profileId", controller.DeletedProfile)
	app.Post("/api/upload_image", controller.UploadImage)

	// add routes in here
}

func (controller *TestRestApiController) TestApiGetList(ctx *fiber.Ctx) error {

	response, err := controller.TestRestApiService.TestApiGetList(ctx)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	return ctx.Status(200).JSON(response)
}

func (controller *TestRestApiController) TestAPiByID(ctx *fiber.Ctx) error {

	profileId := ctx.Params("profileId")
	id, _ := strconv.Atoi(profileId)
	response, err := controller.TestRestApiService.TestApiGetByID(ctx, id)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	return ctx.Status(200).JSON(response)

}

func (controller *TestRestApiController) InsertProfile(ctx *fiber.Ctx) error {
	var request request.ProfileRequest

	err := ctx.BodyParser(&request)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	log.Println(requestData)

	response, err := controller.TestRestApiService.InsertProfile(ctx, &request)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	log.Println(responseData)

	return ctx.Status(200).JSON(response)
}

func (controller *TestRestApiController) Updatedprofile(ctx *fiber.Ctx) error {
	var profileRequest request.ProfileRequest
	err := ctx.BodyParser(&profileRequest)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	profileId := ctx.Params("profileId")
	id, err := strconv.Atoi(profileId)

	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	profileRequest.Id = id

	response, err := controller.TestRestApiService.UpdatedProfile(ctx, &profileRequest)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	log.Println(responseData)

	return ctx.Status(200).JSON(response)
}

func (controller *TestRestApiController) DeletedProfile(ctx *fiber.Ctx) error {

	profileId := ctx.Params("profileId")
	id, _ := strconv.Atoi(profileId)

	response, err := controller.TestRestApiService.DeletedProfile(ctx, id)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	log.Println(responseData)

	return ctx.Status(200).JSON(response)
}

func (controller *TestRestApiController) UploadImage(ctx *fiber.Ctx) error {
	var request request.ProductRequest

	err := ctx.BodyParser(&request)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	log.Println(requestData)

	response, err := controller.TestRestApiService.InsertProduct(ctx, &request)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	log.Println(responseData)

	return ctx.Status(200).JSON(response)
}
