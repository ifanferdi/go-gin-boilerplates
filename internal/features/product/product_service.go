package product

import (
	"strings"

	services "go-gin-boilerplate/internal/features/product/_services"
	"go-gin-boilerplate/internal/helpers"
)

type ProductService interface {
	FindAll(params FindAllProductDTO) (helpers.PaginateResult[Product], error)
	FindByID(id uint) (*Product, error)
	Create(dto CreateProductDTO) (*Product, error)
	Update(dto UpdateProductDTO) (*Product, error)
	Delete(id uint) error
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) FindAll(params FindAllProductDTO) (helpers.PaginateResult[Product], error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	data, err := s.repo.FindAll(params)
	if err != nil {
		return helpers.PaginateResult[Product]{}, err
	}

	total, err := s.repo.Count(params)
	if err != nil {
		return helpers.PaginateResult[Product]{}, err
	}

	return helpers.Paginate(helpers.PaginateParam[Product]{
		Page:  params.Page,
		Limit: params.Limit,
		Total: total,
		Data:  data,
	}), nil
}

func (s *productService) FindByID(id uint) (*Product, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, helpers.NewErrorNotFound("Produk tidak ditemukan")
	}
	return item, nil
}

// Create mendelegasikan logic panjang pembuatan Product + relasi Category ke _services/product_service_create.go
// namun tetap dipanggil oleh controller melalui ProductService (product_service.go).
func (s *productService) Create(dto CreateProductDTO) (*Product, error) {
	return services.CreateProduct[Product](s.repo, services.CreateProductInput{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Cost:        dto.Cost,
		SKU:         dto.SKU,
		Barcode:     dto.Barcode,
		ImagePath:   dto.ImagePath,
		IsActive:    dto.IsActive,
		Stock:       dto.Stock,
		CategoryIDs: dto.CategoryIDs,
	})
}

func (s *productService) Update(dto UpdateProductDTO) (*Product, error) {
	existing, err := s.FindByID(dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.SKU != nil && strings.TrimSpace(*dto.SKU) != "" {
		sku := strings.TrimSpace(*dto.SKU)
		count, err := s.repo.CountBySKU(sku, dto.ID)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, helpers.NewErrorBadRequest("SKU telah digunakan.")
		}
		existing.SKU = &sku
	}

	if dto.Barcode != nil && strings.TrimSpace(*dto.Barcode) != "" {
		barcode := strings.TrimSpace(*dto.Barcode)
		count, err := s.repo.CountByBarcode(barcode, dto.ID)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, helpers.NewErrorBadRequest("Barcode telah digunakan.")
		}
		existing.Barcode = &barcode
	}

	if dto.CategoryIDs != nil && len(*dto.CategoryIDs) > 0 {
		validCount, err := s.repo.CountExistingCategories(*dto.CategoryIDs)
		if err != nil {
			return nil, err
		}
		if int(validCount) != len(*dto.CategoryIDs) {
			return nil, helpers.NewErrorBadRequest("Satu atau lebih categoryId tidak ditemukan.")
		}
	}

	existing.Name = strings.TrimSpace(dto.Name)
	existing.Description = dto.Description
	existing.Price = dto.Price
	existing.Cost = dto.Cost
	existing.ImagePath = dto.ImagePath
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}
	existing.Stock = dto.Stock

	return s.repo.UpdateWithCategories(existing, dto.CategoryIDs)
}

func (s *productService) Delete(id uint) error {
	if _, err := s.FindByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
