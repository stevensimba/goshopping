package entities

type Item struct {
	Product  Product `json:"product"`
	Quantity int64   `json:"quantity"`
}
