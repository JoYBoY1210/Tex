package models

import "time"

type Product struct {
	ID          string   `gorm:"primaryKey" json:"id"`
	Name        string   `gorm:"not null" json:"name"`
	Description string   `gorm:"not null" json:"description"`
	Price       float64  `gorm:"not null" json:"price"`
	CategoryID  string   `gorm:"not null;index" json:"category_id"`
	Category    Category `gorm:"foreignKey:CategoryID;constraint:onUpdate:CASCADE,onDelete:RESTRICT" json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
