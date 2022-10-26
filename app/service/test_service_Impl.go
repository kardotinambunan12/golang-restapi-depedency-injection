package service

import (
	"fmt"
	errorhandler "golang-depedency-injection/app/error_handler"
	"golang-depedency-injection/app/helper"
	"golang-depedency-injection/app/model/request"
	"golang-depedency-injection/app/model/response"
	"golang-depedency-injection/app/repository"
	"regexp"
	"strings"

	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/uuid"
)

func NewTestApiService(testApiRepository *repository.TestRestApiRepository) TestRestApiService {
	return &testRestApiServiceImpl{
		TestRestApiRepository: *testApiRepository,
	}
}

type testRestApiServiceImpl struct {
	TestRestApiRepository repository.TestRestApiRepository
}

func (service *testRestApiServiceImpl) TestApiGetList(ctx *fiber.Ctx) (*response.ProfileResponse, error) {
	profile, err := service.TestRestApiRepository.TestRestApiGetAll(ctx)

	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	requestData := &response.ProfileResponse{
		Code:   profile.Code,
		Status: profile.Status,
		Data:   profile.Data,
	}

	return requestData, nil
}

func (service *testRestApiServiceImpl) TestApiGetByID(ctx *fiber.Ctx, profileId int) (*response.ProfileResponses, error) {

	profile, err := service.TestRestApiRepository.FindById(profileId)
	if err != nil {
		message := err.Error()
		panic(errorhandler.GeneralError{
			Message: message,
		})
	}

	requestData := &response.ProfileResponses{
		Id:      profile.Id,
		Name:    profile.Name,
		Email:   profile.Email,
		Hobby:   profile.Hobby,
		Address: profile.Address,
	}
	fmt.Println("request Data", requestData)

	return requestData, nil
}

func (service *testRestApiServiceImpl) InsertProfile(ctx *fiber.Ctx, params *request.ProfileRequest) (*fiber.Map, error) {
	profileInsert, err := service.TestRestApiRepository.InsertProfile(params)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	fmt.Println(profileInsert)
	result := &fiber.Map{
		"Message": "Success",
	}
	return result, nil

}

func (service *testRestApiServiceImpl) UpdatedProfile(ctx *fiber.Ctx, params *request.ProfileRequest) (*fiber.Map, error) {
	profile, err := service.TestRestApiRepository.UpdatedProfile(ctx, params)

	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	fmt.Println(profile)

	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	result := &fiber.Map{
		"id":      profile.Id,
		"name":    profile.Name,
		"email":   profile.Email,
		"hobby":   profile.Hobby,
		"address": profile.Address,
	}
	return result, nil
}

func (service *testRestApiServiceImpl) DeletedProfile(ctx *fiber.Ctx, profileId int) (*fiber.Map, error) {
	profile, err := service.TestRestApiRepository.DeletedProfile(profileId)

	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	fmt.Println(profile)

	result := &fiber.Map{
		"Message": "Success",
	}
	return result, nil
}

func (service *testRestApiServiceImpl) InsertProduct(ctx *fiber.Ctx, params *request.ProductRequest) (*response.ProductResponse, error) {
	productInsert, err := service.TestRestApiRepository.InsertProduct(params)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	file, errFile := ctx.FormFile("cover")
	if errFile != nil {
		errorhandler.PanicIfNeeded(err)
	}
	uploadDirectory := helper.StorageDirectory()
	ext := strings.ToLower(regexp.MustCompile(".*\\.([^\\.]+)$").ReplaceAllString(file.Filename, "$1"))
	id := uuid.New()
	newFileName := id.String() + "." + ext

	absolutePath := uploadDirectory + "/" + newFileName
	if err := ctx.SaveFile(file, absolutePath); nil != err {
		return nil, nil
	}
	multimedia := response.ProductResponse{}
	// size := float64(file.Size / 1024)
	fname := file.Filename
	if isImage, err := regexp.Match(`\.(png|jpe?g)$`, []byte(strings.ToLower(fname))); nil == err && isImage {
		w, h, err := helper.GetImageScaleSize(absolutePath)
		if nil == err {
			height := float64(h)
			width := float64(w)
			multimedia.DimensionHeight = &height
			multimedia.DimensionWidth = &width
		}
		return nil, nil
	}
	multimedia.Id = id
	multimedia.ImageFile = newFileName
	multimediaURL := viper.GetString("BASE_URL_FILE") + "/files/" + *&multimedia.ImageFile
	multimedia.URL = multimediaURL
	// fileName := file.Filename
	// errSaveFile := ctx.SaveFile(file, fmt.Sprintf("./public/covers/%s", fileName))
	// errorhandler.PanicIfNeeded(errSaveFile)

	fmt.Println(productInsert)
	// result := &fiber.Map{
	// 	"Message": "Success",
	// }
	return &multimedia, nil
}
