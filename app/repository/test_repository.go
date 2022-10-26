package repository

import (
	"golang-depedency-injection/app/model/request"
	"golang-depedency-injection/app/model/response"

	"github.com/gofiber/fiber/v2"
)

type TestRestApiRepository interface {
	TestRestApiGetAll(ctx *fiber.Ctx) (*response.ProfileResponse, error)
	FindById(profileId int) (*response.ProfileResponses, error)
	InsertProfile(params *request.ProfileRequest) (*response.ProfileResponse, error)
	UpdatedProfile(ctx *fiber.Ctx, params *request.ProfileRequest) (*response.ProfileResponses, error)
	DeletedProfile(profileId int) (*response.GeneralResponse, error)

	InsertProduct(params *request.ProductRequest) (*response.ProfileResponse, error)
}
