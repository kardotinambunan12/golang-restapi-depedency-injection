package lib

/*
 * Microservices Upload
 *
 * API version: 1.0.0
 * Contact     : dikhi.martin@tog.co.id
 */

import (
	"bytes"
	"encoding/json"
	"fmt"
	errorhandler "golang-depedency-injection/app/error_handler"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type UploadFileMultipleStruct struct {
	Filename string `json:"Filename"`
	Header   struct {
		ContentDisposition []string `json:"Content-Disposition"`
		ContentType        []string `json:"Content-Type"`
	} `json:"Header"`
	Size int `json:"Size"`
}

type ResponseUpload struct {
	Filename     string         `json:"file_name"`
	RelativePath string         `json:"relative_path"`
	Status       ResponseStatus `json:"status"`
}

type ResponseStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CheckConnection(url_api string) bool {
	var requestBody bytes.Buffer
	req, err := http.NewRequest("POST", url_api, &requestBody)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		return false
	}

	// Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		return false
	}
	var result map[string]map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)

	if result["status"]["code"] != "200" {
		return false
	}
	return true
}

// ==> Receipt By Form
func UploadFile(c *fiber.Ctx, url_api, path, file_name, field_image string) ResponseUpload {
	form, err := c.MultipartForm()
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	// check_file_empty
	files := form.File[field_image]
	if files == nil {
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    204,
				Message: "no_content_found",
			},
		}
		return response
	}

	// check_connection_microservices
	microservices_connect := CheckConnection(GetEnv("HOST_CDN_IMAGE") + GetEnv("URI_CDN_CHECK_CONNECTION"))
	if microservices_connect == false {
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    503,
				Message: "service_unavailable",
			},
		}
		return response
	}

	file, _ := c.FormFile(field_image)
	file_image, _ := file.Open()
	defer file_image.Close()

	if file_name == "" {
		file_name = file.Filename
	} else {
		file_name = file_name + "_" + file.Filename
	}

	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)
	fileWriter, err := multiPartWriter.CreateFormFile(field_image, file_name)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	_, err = io.Copy(fileWriter, file_image)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	fieldPath, err := multiPartWriter.CreateFormField("path")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	_, err = fieldPath.Write([]byte(path))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	fieldImage, err := multiPartWriter.CreateFormField("field_image")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	_, err = fieldImage.Write([]byte(field_image))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	multiPartWriter.Close()

	req, err := http.NewRequest("POST", url_api, &requestBody)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)
	file_name = fmt.Sprintf("%v", result["data"])

	response_data := ResponseUpload{
		Filename:     file_name,
		RelativePath: path + "/" + file_name,
		Status: ResponseStatus{
			Code:    200,
			Message: "OK",
		},
	}

	return response_data
}

func ReUploadFile(c *fiber.Ctx, url_api, path, old_file, file_name, field_image string) ResponseUpload {
	form, err := c.MultipartForm()
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	// check_file_empty
	files := form.File[field_image]
	if files == nil {
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    204,
				Message: "no_content_found",
			},
		}
		return response
	}

	// check_connection_microservices
	microservices_connect := CheckConnection(GetEnv("HOST_CDN_IMAGE") + GetEnv("URI_CDN_CHECK_CONNECTION"))
	if microservices_connect == false {
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    503,
				Message: "service_unavailable",
			},
		}
		return response
	}

	file, _ := c.FormFile(field_image)
	file_image, _ := file.Open()
	defer file_image.Close()

	if file_name == "" {
		file_name = file.Filename
	} else {
		file_name = file_name + "_" + file.Filename
	}

	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)

	fileWriter, err := multiPartWriter.CreateFormFile(field_image, file_name)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	_, err = io.Copy(fileWriter, file_image)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	fieldPath, err := multiPartWriter.CreateFormField("path")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	_, err = fieldPath.Write([]byte(path))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	if old_file != "" {
		fieldOldFile, err := multiPartWriter.CreateFormField("old_file")
		if err != nil {
			errorhandler.PanicIfNeeded(err)
			response := ResponseUpload{
				Status: ResponseStatus{
					Code:    500,
					Message: "internal_server_error",
				},
			}
			return response
		}
		_, err = fieldOldFile.Write([]byte(old_file))
		if err != nil {
			errorhandler.PanicIfNeeded(err)
			response := ResponseUpload{
				Status: ResponseStatus{
					Code:    500,
					Message: "internal_server_error",
				},
			}
			return response
		}
	}

	fieldImage, err := multiPartWriter.CreateFormField("field_image")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	_, err = fieldImage.Write([]byte(field_image))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	multiPartWriter.Close()

	req, err := http.NewRequest("POST", url_api, &requestBody)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	// result
	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)

	file_name = fmt.Sprintf("%v", result["data"])
	response_data := ResponseUpload{
		Filename:     file_name,
		RelativePath: path + "/" + file_name,
		Status: ResponseStatus{
			Code:    200,
			Message: "OK",
		},
	}

	return response_data
}

func GetEnv(val string) string {
	e := godotenv.Load(".env")
	if e != nil {
		panic(e)
	}
	return os.Getenv(val)
}

func DeleteFile(c *fiber.Ctx, url_api, path, old_file string) ResponseUpload {

	// check_connection_microservices
	microservices_connect := CheckConnection(GetEnv("HOST_CDN_IMAGE") + GetEnv("URI_CDN_CHECK_CONNECTION"))
	if microservices_connect == false {
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    503,
				Message: "service_unavailable",
			},
		}
		return response
	}

	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)

	fieldPath, err := multiPartWriter.CreateFormField("path")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	_, err = fieldPath.Write([]byte(path))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	fieldOldFile, err := multiPartWriter.CreateFormField("old_file")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	_, err = fieldOldFile.Write([]byte(old_file))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	multiPartWriter.Close()

	req, err := http.NewRequest("POST", url_api, &requestBody)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	response_data := ResponseUpload{
		Status: ResponseStatus{
			Code:    200,
			Message: "OK",
		},
	}
	return response_data
}

// ==> Receipt Base64
func UploadFileBase64(base_64, image_name, url_api, path string) ResponseUpload {
	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)

	// check_connection_microservices
	microservices_connect := CheckConnection(GetEnv("HOST_CDN_IMAGE") + GetEnv("URI_CDN_CHECK_CONNECTION"))
	if microservices_connect == false {
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    503,
				Message: "service_unavailable",
			},
		}
		return response
	}

	// Image Form
	fieldImageName, err := multiPartWriter.CreateFormField("image_name")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	_, err = fieldImageName.Write([]byte(image_name))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	fieldBase64, err := multiPartWriter.CreateFormField("base_64")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	_, err = fieldBase64.Write([]byte(base_64))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	// Additional Form
	fieldPath, err := multiPartWriter.CreateFormField("path")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	_, err = fieldPath.Write([]byte(path))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	req, err := http.NewRequest("POST", url_api, &requestBody)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)
	file_name := fmt.Sprintf("%v", result["data"])

	response_data := ResponseUpload{
		Filename:     file_name,
		RelativePath: path + "/" + file_name,
		Status: ResponseStatus{
			Code:    200,
			Message: "OK",
		},
	}

	return response_data
}

func ReUploadFileBase64(base_64 []byte, image_name, url_api, path, old_file string) ResponseUpload {
	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)

	// check_connection_microservices
	microservices_connect := CheckConnection(GetEnv("HOST_CDN_IMAGE") + GetEnv("URI_CDN_CHECK_CONNECTION"))
	if microservices_connect == false {
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    503,
				Message: "service_unavailable",
			},
		}
		return response
	}

	// Image Form
	fieldImageName, err := multiPartWriter.CreateFormField("image_name")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	_, err = fieldImageName.Write([]byte(image_name))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	fieldBase64, err := multiPartWriter.CreateFormField("base_64")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	_, err = fieldBase64.Write([]byte(base_64))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	// Additional Form
	fieldPath, err := multiPartWriter.CreateFormField("path")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}
	_, err = fieldPath.Write([]byte(path))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	if old_file != "" {
		fieldOldFile, err := multiPartWriter.CreateFormField("old_file")
		if err != nil {
			errorhandler.PanicIfNeeded(err)
			response := ResponseUpload{
				Status: ResponseStatus{
					Code:    500,
					Message: "internal_server_error",
				},
			}
			return response
		}
		_, err = fieldOldFile.Write([]byte(old_file))
		if err != nil {
			errorhandler.PanicIfNeeded(err)
			response := ResponseUpload{
				Status: ResponseStatus{
					Code:    500,
					Message: "internal_server_error",
				},
			}
			return response
		}
	}

	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	req, err := http.NewRequest("POST", url_api, &requestBody)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	// Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		response := ResponseUpload{
			Status: ResponseStatus{
				Code:    500,
				Message: "internal_server_error",
			},
		}
		return response
	}

	// result
	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)

	file_name := fmt.Sprintf("%v", result["data"])
	response_data := ResponseUpload{
		Filename:     file_name,
		RelativePath: path + "/" + file_name,
		Status: ResponseStatus{
			Code:    200,
			Message: "OK",
		},
	}

	return response_data
}

func UploadFileMultiple(c *fiber.Ctx, url_api, path, field_image string) []UploadFileMultipleStruct {
	form, err := c.MultipartForm()
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		return nil
	}
	files := form.File[field_image]
	if files == nil {
		errorhandler.PanicIfNeeded(err)
		return nil
	}

	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)

	for key, value := range files {
		file, _ := files[key].Open()
		defer file.Close()

		// Initialize the file field
		fileWriter, err := multiPartWriter.CreateFormFile(field_image, value.Filename)
		if err != nil {
			errorhandler.PanicIfNeeded(err)
		}

		// Copy the actual file content to the field field's writer
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			errorhandler.PanicIfNeeded(err)
		}
	}

	fieldPath, err := multiPartWriter.CreateFormField("path")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	_, err = fieldPath.Write([]byte(path))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	fieldImage, err := multiPartWriter.CreateFormField("field_image")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	_, err = fieldImage.Write([]byte(field_image))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	multiPartWriter.Close()

	req, err := http.NewRequest("POST", url_api, &requestBody)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	var json_data map[string]interface{}
	json.NewDecoder(response.Body).Decode(&json_data)
	jsonString, _ := json.Marshal(json_data["data"])

	jsonByteArray := []byte(string(jsonString))

	var result []UploadFileMultipleStruct
	err = json.Unmarshal(jsonByteArray, &result)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	return result
}

func ReUploadFileMultiple(c *fiber.Ctx, url_api, path, field_image string) []UploadFileMultipleStruct {

	form, err := c.MultipartForm()
	if err != nil {
		errorhandler.PanicIfNeeded(err)
		return nil
	}
	files := form.File[field_image]
	if files == nil {
		errorhandler.PanicIfNeeded(err)
		return nil
	}

	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)

	// multiple_upload
	for key, value := range files {
		file, _ := files[key].Open()
		defer file.Close()

		// Initialize the file field
		fileWriter, err := multiPartWriter.CreateFormFile(field_image, value.Filename)
		if err != nil {
			errorhandler.PanicIfNeeded(err)
		}

		// Copy the actual file content to the field field's writer
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			errorhandler.PanicIfNeeded(err)
		}
	}

	fieldPath, err := multiPartWriter.CreateFormField("path")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	_, err = fieldPath.Write([]byte(path))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	fieldImage, err := multiPartWriter.CreateFormField("field_image")
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	_, err = fieldImage.Write([]byte(field_image))
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}
	multiPartWriter.Close()

	req, err := http.NewRequest("POST", url_api, &requestBody)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	var json_data map[string]interface{}
	json.NewDecoder(response.Body).Decode(&json_data)
	jsonString, _ := json.Marshal(json_data["data"])

	jsonByteArray := []byte(string(jsonString))

	var result []UploadFileMultipleStruct
	err = json.Unmarshal(jsonByteArray, &result)
	if err != nil {
		errorhandler.PanicIfNeeded(err)
	}

	return result
}
