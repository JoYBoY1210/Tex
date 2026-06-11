package db

import (
	"log"

	"github.com/joyboy1210/tex/internal/models"
	"gorm.io/gorm"
)

func SeedDb(db *gorm.DB) error {
	var count int64

	err := db.Model(&models.Category{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		log.Printf("Db is already seeded")
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		log.Printf("Seeding database with initial data")

		laLiga := models.Category{
			Name: "La liga",
		}
		bundesliga := models.Category{
			Name: "Bundesliga",
		}
		premierLeague := models.Category{
			Name: "Premier league",
		}
		accessories := models.Category{
			Name: "Accessories",
		}

		if err := tx.Create(&laLiga).Error; err != nil {
			return err
		}

		if err := tx.Create(&bundesliga).Error; err != nil {
			return err
		}

		if err := tx.Create(&premierLeague).Error; err != nil {
			return err
		}

		if err := tx.Create(&accessories).Error; err != nil {
			return err
		}

		products := []models.Product{
			{
				CategoryID:  laLiga.ID,
				Name:        "Real Madrid Home 24/25",
				Description: "The classic iconic white home kit. Hala Madrid!",
				Price:       89.99,
			},
			{
				CategoryID:  laLiga.ID,
				Name:        "Barcelona Away 24/25",
				Description: "Blacked out away kit with red/blue trims.",
				Price:       85.00,
			},
			{
				CategoryID:  premierLeague.ID,
				Name:        "Arsenal Home 24/25",
				Description: "Red with white sleeves. Classic Gunners.",
				Price:       90.00,
			},
			{
				CategoryID:  premierLeague.ID,
				Name:        "Man City Home 24/25",
				Description: "Sky blue home kit. Champions edition.",
				Price:       95.00,
			},
			{
				CategoryID:  accessories.ID,
				Name:        "UCL Official Match Ball",
				Description: "Champions League 24/25 Official Match Ball.",
				Price:       45.00,
			},
			{
				CategoryID:  bundesliga.ID,
				Name:        "Bayern Munich Home 24/25",
				Description: "The classic iconic red home kit. FC Bayern München!",
				Price:       90.00,
			},
			{
				CategoryID:  bundesliga.ID,
				Name:        "Borussia Dortmund Home 24/25",
				Description: "The classic iconic yellow home kit. BVB!",
				Price:       85.00,
			},
		}

		if err := tx.Create(&products).Error; err != nil {
			return err
		}

		log.Printf("Database seeded successfully")
		return nil
	})
}
