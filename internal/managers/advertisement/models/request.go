package models

type CreateAdvertisementRequest struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Price       float32     `json:"price"`
	Quantity    int32       `json:"quantity"`
	Photos      []PhotoFile `json:"photos"`
}

type PhotoFile struct {
	Name    string
	Content []byte
}
