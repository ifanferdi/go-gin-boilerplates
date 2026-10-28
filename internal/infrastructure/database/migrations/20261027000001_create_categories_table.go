package migrations

import (
	"go-gin-boilerplate/internal/features/category"

	"gorm.io/gorm"
)

func init() {
	RegisterMigration(
		"20261027000001_create_categories_table",
		func(db *gorm.DB) error {
			return db.AutoMigrate(&category.Category{})
		},
		func(db *gorm.DB) error {
			return db.Migrator().DropTable(&category.Category{})
		},
	)
}
