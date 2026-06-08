package models

import "time"

type User struct {
	PhoneNumber  string     `gorm:"primaryKey" json:"phone_number"`
	CurrentState string     `gorm:"notNull" json:"current_state"`
	Cart         []CartItem `gorm:"serializer:json" json:"cart"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func UpdateUserState(phone, newState string) error {
	result := DB.Model(&User{}).Where("phone_number = ?", phone).Update("current_state", newState)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUser(phone string) (*User, error) {
	var user User
	result := DB.First(&user, "phone_number = ?", phone)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func CreateUser(phone string,state string) (*User, error) {
	user := User{
		PhoneNumber:  phone,
		CurrentState: state,
		Cart:         []CartItem{},
	}
	result := DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
