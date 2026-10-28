package services

import (
	"strings"

	"go-gin-boilerplate/internal/helpers"
)

type CreateProductInput struct {
	Name        string
	Description *string
	Price       int
	Cost        *int
	SKU         *string
	Barcode     *string
	ImagePath   *string
	IsActive    *bool
	Stock       int
	CategoryIDs []uint
}

type ProductCreatorRepository[T any] interface {
	CountBySKU(sku string, excludeID uint) (int64, error)
	CountByBarcode(barcode string, excludeID uint) (int64, error)
	CountExistingCategories(ids []uint) (int64, error)
	CreateWithCategories(input CreateProductInput) (*T, error)
}

// CreateProduct menangani business logic panjang pembuatan Product + validasi SKU, Barcode, dan relasi Category.
// Dipanggil oleh controller melalui ProductService.Create() di product_service.go.
func CreateProduct[T any](repo ProductCreatorRepository[T], input CreateProductInput) (*T, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return nil, helpers.NewErrorBadRequest("Column name: Name cannot be empty.")
	}

	if input.SKU != nil && strings.TrimSpace(*input.SKU) != "" {
		sku := strings.TrimSpace(*input.SKU)
		input.SKU = &sku
		count, err := repo.CountBySKU(sku, 0)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, helpers.NewErrorBadRequest("SKU telah digunakan.")
		}
	}

	if input.Barcode != nil && strings.TrimSpace(*input.Barcode) != "" {
		barcode := strings.TrimSpace(*input.Barcode)
		input.Barcode = &barcode
		count, err := repo.CountByBarcode(barcode, 0)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, helpers.NewErrorBadRequest("Barcode telah digunakan.")
		}
	}

	if len(input.CategoryIDs) > 0 {
		uniqueIDs := uniqueUintSlice(input.CategoryIDs)
		input.CategoryIDs = uniqueIDs
		validCount, err := repo.CountExistingCategories(uniqueIDs)
		if err != nil {
			return nil, err
		}
		if int(validCount) != len(uniqueIDs) {
			return nil, helpers.NewErrorBadRequest("Satu atau lebih categoryId tidak ditemukan.")
		}
	}

	return repo.CreateWithCategories(input)
}

func uniqueUintSlice(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			out = append(out, id)
		}
	}
	return out
}
