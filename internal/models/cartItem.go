package models


type CartItem struct {
	ProductID uint `json:"product_id"`
	Quantity int `json:"quantity"`
}