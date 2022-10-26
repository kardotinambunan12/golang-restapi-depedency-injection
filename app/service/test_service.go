package service

import (
	"golang-depedency-injection/app/model/request"
	"golang-depedency-injection/app/model/response"

	"github.com/gofiber/fiber/v2"
)

type TestRestApiService interface {
	TestApiGetList(ctx *fiber.Ctx) (*response.ProfileResponse, error)
	TestApiGetByID(ctx *fiber.Ctx, profileId int) (*response.ProfileResponses, error)
	InsertProfile(ctx *fiber.Ctx, params *request.ProfileRequest) (*fiber.Map, error)
	UpdatedProfile(ctx *fiber.Ctx, params *request.ProfileRequest) (*fiber.Map, error)
	DeletedProfile(ctx *fiber.Ctx, profileId int) (*fiber.Map, error)
	InsertProduct(ctx *fiber.Ctx, params *request.ProductRequest) (*response.ProductResponse, error)
}
