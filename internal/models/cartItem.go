package models


type CartItem struct {
	ProductID string `json:"product_id"`
	Quantity int `json:"quantity"`
}