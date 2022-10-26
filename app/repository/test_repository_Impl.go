package repository

import (
	"errors"
	"fmt"
	"golang-depedency-injection/app/config"
	errorhandler "golang-depedency-injection/app/error_handler"
	"golang-depedency-injection/app/model/request"
	"golang-depedency-injection/app/model/response"

	"github.com/gofiber/fiber/v2"
)

func NewTestApiRepository() TestRestApiRepository {
	return &testApiRepositoryImpl{
		// Configuration: *configuration,
	}
}

type testApiRepositoryImpl struct {
	// Configuration config.Config
}

func (repository *testApiRepositoryImpl) TestRestApiGetAll(ctx *fiber.Ctx) (*response.ProfileResponse, error) {
	db := config.NewDB()
	sql := `select id, name, email, hobby, address from profile`
	rows, err := db.Query(sql)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	profiles := make([]response.ProfileResponses, 0)
	for rows.Next() {
		profile := response.ProfileResponses{}
		err := rows.Scan(&profile.Id, &profile.Name, &profile.Email, &profile.Hobby, &profile.Address)
		if err != nil {
			errorhandler.PanicIfNeeded(err)
		}
		profiles = append(profiles, profile)
	}

	data := &response.ProfileResponse{
		Code:   200,
		Status: "Success",
		Data:   profiles,
	}
	return data, nil

}

func (repository *testApiRepositoryImpl) FindById(profileId int) (*response.ProfileResponses, error) {

	db := config.NewDB()

	SQL := "select id, name, email, hobby, address from profile where id = ?"
	rows, err := db.Query(SQL, profileId)
	fmt.Println("profile row", profileId)
	errorhandler.PanicIfNeeded(err)
	defer rows.Close()

	profile := &response.ProfileResponses{}
	if rows.Next() {
		err := rows.Scan(&profile.Id, &profile.Name, &profile.Email, &profile.Hobby, &profile.Address)
		errorhandler.PanicIfNeeded(err)
		fmt.Println("profile res", profile)
		return profile, nil
	} else {
		return profile, errors.New("profile is not found")
	}
}

func (repository *testApiRepositoryImpl) InsertProfile(params *request.ProfileRequest) (*response.ProfileResponse, error) {
	db := config.NewDB()
	defer db.Close()
	sql := `insert into profile(name, email, hobby, address) values(?, ?, ?, ?)`
	result, err := db.Exec(sql, params.Name, params.Email, params.Hobby, params.Address)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	// id, err := result.LastInsertId()
	// errorhandler.PanicIfNeeded(err)
	fmt.Println(result)

	resData := &response.ProfileResponse{
		Code:   200,
		Status: "Success",
	}

	return resData, nil
}

func (repository *testApiRepositoryImpl) UpdatedProfile(ctx *fiber.Ctx, params *request.ProfileRequest) (*response.ProfileResponses, error) {
	db := config.NewDB()
	defer db.Close()

	sql := "update profile set name = ?, email=?, hobby=?, address=? where id = ?"
	result, err := db.Exec(sql, params.Name, params.Email, params.Hobby, params.Address, params.Id)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	// id, err := result.LastInsertId()
	// errorhandler.PanicIfNeeded(err)
	fmt.Println(result)

	resData := &response.ProfileResponses{
		Name:    params.Name,
		Email:   params.Email,
		Hobby:   params.Hobby,
		Address: params.Address,
	}

	return resData, nil
}

func (repository *testApiRepositoryImpl) DeletedProfile(profileId int) (*response.GeneralResponse, error) {
	db := config.NewDB()
	defer db.Close()

	SQL := "delete from profile where id = ?"
	_, err := db.Exec(SQL, profileId)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	resData := &response.GeneralResponse{
		StatusCode: 200,
		Message:    "Success",
	}

	return resData, nil
}

func (repository *testApiRepositoryImpl) InsertProduct(params *request.ProductRequest) (*response.ProfileResponse, error) {
	db := config.NewDB()
	defer db.Close()
	sql := `insert into product(name, imagefile, description, harga) values(?, ?, ?, ?)`
	result, err := db.Exec(sql, params.Name, params.ImageFile, params.Description, params.Harga)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	fmt.Println(result)

	resData := &response.ProfileResponse{
		Code:   200,
		Status: "Success",
	}

	return resData, nil
}
