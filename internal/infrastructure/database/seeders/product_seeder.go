package seeders

import (
	"log"

	"go-gin-boilerplate/internal/features/category"
	"go-gin-boilerplate/internal/features/product"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedProducts(db *gorm.DB) error {
	strPtr := func(s string) *string { return &s }
	intPtr := func(i int) *int { return &i }

	var categories []category.Category
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	catMap := make(map[string]uint)
	for _, c := range categories {
		catMap[c.Name] = c.ID
	}

	initialProducts := []product.Product{
		{
			Name:        "Iced Americano",
			Description: strPtr("Double shot espresso over ice"),
			Price:       25000,
			Cost:        intPtr(10000),
			SKU:         strPtr("SKU-COFFEE-001"),
			Barcode:     strPtr("899000000001"),
			IsActive:    true,
			Stock:       100,
		},
		{
			Name:        "Matcha Latte",
			Description: strPtr("Premium Uji matcha with fresh milk"),
			Price:       32000,
			Cost:        intPtr(14000),
			SKU:         strPtr("SKU-BEV-001"),
			Barcode:     strPtr("899000000002"),
			IsActive:    true,
			Stock:       80,
		},
		{
			Name:        "Butter Croissant",
			Description: strPtr("Flaky French butter croissant"),
			Price:       28000,
			Cost:        intPtr(12000),
			SKU:         strPtr("SKU-SNK-001"),
			Barcode:     strPtr("899000000003"),
			IsActive:    true,
			Stock:       50,
		},
	}

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "sku"}},
		DoNothing: true,
	}).Create(&initialProducts).Error; err != nil {
		return err
	}

	// Seed pivot table product_has_categories
	var seededProducts []product.Product
	if err := db.Find(&seededProducts).Error; err != nil {
		return err
	}

	var pivots []product.ProductHasCategory
	for _, p := range seededProducts {
		if p.SKU == nil {
			continue
		}
		switch *p.SKU {
		case "SKU-COFFEE-001":
			if id, ok := catMap["Coffee"]; ok {
				pivots = append(pivots, product.ProductHasCategory{ProductID: p.ID, CategoryID: id})
			}
			if id, ok := catMap["Beverages"]; ok {
				pivots = append(pivots, product.ProductHasCategory{ProductID: p.ID, CategoryID: id})
			}
		case "SKU-BEV-001":
			if id, ok := catMap["Beverages"]; ok {
				pivots = append(pivots, product.ProductHasCategory{ProductID: p.ID, CategoryID: id})
			}
		case "SKU-SNK-001":
			if id, ok := catMap["Snacks"]; ok {
				pivots = append(pivots, product.ProductHasCategory{ProductID: p.ID, CategoryID: id})
			}
			if id, ok := catMap["Food"]; ok {
				pivots = append(pivots, product.ProductHasCategory{ProductID: p.ID, CategoryID: id})
			}
		}
	}

	if len(pivots) > 0 {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&pivots).Error; err != nil {
			return err
		}
	}

	log.Println("✅ Seed Products & ProductHasCategories")
	return nil
}
