package entities

type Product struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
	Photo    string  `json:"photo"`
}
