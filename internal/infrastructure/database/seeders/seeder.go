package seeders

import (
	"log"

	"gorm.io/gorm"
)

func RunAllSeeders(db *gorm.DB) error {
	log.Println("🌱 Running database seeders...")

	if err := SeedCategories(db); err != nil {
		return err
	}

	if err := SeedProducts(db); err != nil {
		return err
	}

	log.Println("✅ All seeders completed")
	return nil
}
