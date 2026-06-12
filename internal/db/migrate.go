package db

import (
	"fmt"

	"github.com/joyboy1210/tex/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Order{})
	if err != nil {
		return err
	}
	fmt.Println("Database migrated successfully")
	return nil
}
