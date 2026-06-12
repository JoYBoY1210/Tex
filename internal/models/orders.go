package models

type Order struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	PhoneNumber string     `gorm:"not null" json:"phone_number"`
	Products    []CartItem `gorm:"serializer:json" json:"products"`
	TotalPrice  float64    `gorm:"not null" json:"total_price"`
	Status      string     `gorm:"not null" json:"status"`
}

func CreateOrder(phone string, products []CartItem, totalPrice float64) (*Order, error) {
	order := Order{
		PhoneNumber: phone,
		Products:    products,
		TotalPrice:  totalPrice,
		Status:      "Pending",
	}
	result := DB.Create(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func GetOrdersByPhone(phone string) ([]Order, error) {
	var orders []Order
	result := DB.Where("phone_number = ?", phone).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func UpdateOrderStatus(orderId uint, newStatus string) error {
	result := DB.Model(&Order{}).Where("id = ?", orderId).Update("status", newStatus)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetOrderStatus(orderId uint) (string, error) {
	var order Order
	result := DB.First(&order, "id = ?", orderId)
	if result.Error != nil {
		return "", result.Error
	}
	return order.Status, nil
}

func GetActiveOrders(phone string) ([]Order, error) {
	var orders []Order
	result := DB.Where("phone_number = ? AND status != ?", phone, "Delivered").Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}
