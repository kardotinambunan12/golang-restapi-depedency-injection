package response

type GeneralResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
