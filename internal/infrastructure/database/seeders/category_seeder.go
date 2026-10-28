package seeders

import (
	"log"

	"go-gin-boilerplate/internal/features/category"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedCategories(db *gorm.DB) error {
	strPtr := func(s string) *string { return &s }

	initialCategories := []category.Category{
		{
			Name:        "Coffee",
			Description: strPtr("Espresso based and manual brew coffee"),
		},
		{
			Name:        "Beverages",
			Description: strPtr("Non-coffee drinks, tea, and refreshers"),
		},
		{
			Name:        "Food",
			Description: strPtr("Main courses and heavy meals"),
		},
		{
			Name:        "Snacks",
			Description: strPtr("Light bites and pastries"),
		},
	}

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).Create(&initialCategories).Error; err != nil {
		return err
	}

	log.Println("✅ Seed Categories")
	return nil
}
