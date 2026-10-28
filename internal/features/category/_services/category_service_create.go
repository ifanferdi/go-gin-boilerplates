package services

import (
	"strings"

	"go-gin-boilerplate/internal/helpers"
)

type CategoryCreatorRepository[T any] interface {
	CountByName(name string, excludeID uint) (int64, error)
	CreateRecord(name string, description *string) (*T, error)
}

// CreateCategory menangani business logic pembuatan Category dan dipanggil melalui CategoryService.Create().
func CreateCategory[T any](repo CategoryCreatorRepository[T], name string, description *string) (*T, error) {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return nil, helpers.NewErrorBadRequest("Column name: Name cannot be empty.")
	}

	count, err := repo.CountByName(cleanName, 0)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, helpers.NewErrorBadRequest("Nama kategori telah digunakan.")
	}

	if description != nil {
		trimmed := strings.TrimSpace(*description)
		description = &trimmed
	}

	return repo.CreateRecord(cleanName, description)
}
