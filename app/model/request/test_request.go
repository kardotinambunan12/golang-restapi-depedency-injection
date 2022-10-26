package request

type ProfileRequest struct {
	// Id      uuid.UUID `json:"id,omitempty" gorm:"primaryKey;unique;type:varchar(36);not null" format:"uuid"`
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Hobby   string `json:"hobby"`
	Address string `json:"address"`
}

type ProfileByIDRequest struct {
	Id int `json:"id"`
}

type ProductRequest struct {
	// Id      uuid.UUID `json:"id,omitempty" gorm:"primaryKey;unique;type:varchar(36);not null" format:"uuid"`
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	ImageFile   string  `json:"imageFile"`
	URL         string  `json:"url"`
	Description string  `json:"description"`
	Harga       float64 `json:"harga"`
}
