package models

import "time"

type User struct {
	PhoneNumber  string     `gorm:"primaryKey" json:"phone_number"`
	CurrentState string     `gorm:"notNull" json:"current_state"`
	Cart         []CartItem `gorm:"serializer:json" json:"cart"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
