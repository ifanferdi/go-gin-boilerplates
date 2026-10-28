package migrations

import (
	"go-gin-boilerplate/internal/features/product"

	"gorm.io/gorm"
)

func init() {
	RegisterMigration(
		"20261027000002_create_products_and_pivot_table",
		func(db *gorm.DB) error {
			if err := db.SetupJoinTable(&product.Product{}, "Categories", &product.ProductHasCategory{}); err != nil {
				return err
			}
			return db.AutoMigrate(&product.Product{}, &product.ProductHasCategory{})
		},
		func(db *gorm.DB) error {
			if err := db.Migrator().DropTable(&product.ProductHasCategory{}); err != nil {
				return err
			}
			return db.Migrator().DropTable(&product.Product{})
		},
	)
}
