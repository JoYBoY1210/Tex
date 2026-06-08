package models

import "gorm.io/gorm"


var DB *gorm.DB

func InitDB(database *gorm.DB){
	DB=database
}