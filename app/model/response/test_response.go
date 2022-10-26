package response

import "github.com/google/uuid"

type ProfileResponse struct {
	Code   int                `json:"code"`
	Status string             `json:"status"`
	Data   []ProfileResponses `json:"data"`
}

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ProfileResponses struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Hobby   string `json:"hobby"`
	Address string `json:"address"`
}

type ProductResponse struct {
	Id uuid.UUID `json:"id,omitempty" gorm:"primaryKey;unique;type:varchar(36);not null" format:"uuid"`
	// Id              int      `json:"id"`
	Name            string   `json:"name"`
	ImageFile       string   `json:"imageFile"`
	URL             string   `json:"url"`
	Description     string   `json:"description"`
	Harga           float64  `json:"harga"`
	DimensionWidth  *float64 `json:"dimension_width"` // Dimension Width
	DimensionHeight *float64 `json:"dimension_height"`
}
