package product

import (
	"errors"

	"go-gin-boilerplate/internal/features/category"
	services "go-gin-boilerplate/internal/features/product/_services"
	"go-gin-boilerplate/internal/helpers"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(params FindAllProductDTO) ([]Product, error)
	Count(params FindAllProductDTO) (int64, error)
	FindByID(id uint) (*Product, error)
	CountBySKU(sku string, excludeID uint) (int64, error)
	CountByBarcode(barcode string, excludeID uint) (int64, error)
	CountExistingCategories(ids []uint) (int64, error)
	CreateWithCategories(input services.CreateProductInput) (*Product, error)
	UpdateWithCategories(product *Product, categoryIDs *[]uint) (*Product, error)
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) applyFilters(query *gorm.DB, params FindAllProductDTO) *gorm.DB {
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("products.name ILIKE ? OR products.description ILIKE ? OR products.sku ILIKE ?", searchPattern, searchPattern, searchPattern)
	}
	if params.SKU != "" {
		query = query.Where("products.sku = ?", params.SKU)
	}
	if params.Barcode != "" {
		query = query.Where("products.barcode = ?", params.Barcode)
	}
	if params.IsActive != nil {
		query = query.Where("products.is_active = ?", *params.IsActive)
	}
	if len(params.IDs) > 0 {
		query = query.Where("products.id IN ?", params.IDs)
	}
	if params.CategoryID > 0 {
		query = query.Joins("JOIN product_has_categories phc ON phc.product_id = products.id").
			Where("phc.category_id = ?", params.CategoryID)
	}
	return query
}

func (r *productRepository) FindAll(params FindAllProductDTO) ([]Product, error) {
	var products []Product

	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit == 0 {
		limit = 10
	}

	query := r.applyFilters(r.db.Model(&Product{}).Preload("Categories"), params)
	orderClause := helpers.ParseOrderBy(params.OrderBy, AllowedOrderColumns, "id DESC")
	query = query.Order("products." + orderClause)

	if limit != -1 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) Count(params FindAllProductDTO) (int64, error) {
	var total int64
	query := r.applyFilters(r.db.Model(&Product{}), params)
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *productRepository) FindByID(id uint) (*Product, error) {
	var item Product
	if err := r.db.Preload("Categories").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *productRepository) CountBySKU(sku string, excludeID uint) (int64, error) {
	var count int64
	query := r.db.Model(&Product{}).Where("sku = ?", sku)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productRepository) CountByBarcode(barcode string, excludeID uint) (int64, error) {
	var count int64
	query := r.db.Model(&Product{}).Where("barcode = ?", barcode)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productRepository) CountExistingCategories(ids []uint) (int64, error) {
	var count int64
	if err := r.db.Model(&category.Category{}).Where("id IN ?", ids).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productRepository) CreateWithCategories(input services.CreateProductInput) (*Product, error) {
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	item := &Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Cost:        input.Cost,
		SKU:         input.SKU,
		Barcode:     input.Barcode,
		ImagePath:   input.ImagePath,
		IsActive:    isActive,
		Stock:       input.Stock,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(item).Error; err != nil {
			return err
		}

		if len(input.CategoryIDs) > 0 {
			pivots := make([]ProductHasCategory, 0, len(input.CategoryIDs))
			for _, catID := range input.CategoryIDs {
				pivots = append(pivots, ProductHasCategory{
					ProductID:  item.ID,
					CategoryID: catID,
				})
			}
			if err := tx.Create(&pivots).Error; err != nil {
				return err
			}
		}

		return tx.Preload("Categories").First(item, item.ID).Error
	})
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *productRepository) UpdateWithCategories(item *Product, categoryIDs *[]uint) (*Product, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(item).Error; err != nil {
			return err
		}

		if categoryIDs != nil {
			if err := tx.Where("product_id = ?", item.ID).Delete(&ProductHasCategory{}).Error; err != nil {
				return err
			}
			if len(*categoryIDs) > 0 {
				pivots := make([]ProductHasCategory, 0, len(*categoryIDs))
				for _, catID := range *categoryIDs {
					pivots = append(pivots, ProductHasCategory{
						ProductID:  item.ID,
						CategoryID: catID,
					})
				}
				if err := tx.Create(&pivots).Error; err != nil {
					return err
				}
			}
		}

		return tx.Preload("Categories").First(item, item.ID).Error
	})
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&Product{}, id).Error
}
