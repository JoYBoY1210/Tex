package models

import "time"

type Category struct {
	ID string `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null;uniqueIndex" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}