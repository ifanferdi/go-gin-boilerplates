package category

import (
	"errors"

	"go-gin-boilerplate/internal/helpers"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(params FindAllCategoryDTO) ([]Category, error)
	Count(params FindAllCategoryDTO) (int64, error)
	FindByID(id uint) (*Category, error)
	FindByIDs(ids []uint) ([]Category, error)
	CountByName(name string, excludeID uint) (int64, error)
	CreateRecord(name string, description *string) (*Category, error)
	Create(category *Category) error
	Update(category *Category) error
	Delete(id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) applyFilters(query *gorm.DB, params FindAllCategoryDTO) *gorm.DB {
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("product_categories.name ILIKE ? OR product_categories.description ILIKE ?", searchPattern, searchPattern)
	}
	if params.Name != "" {
		query = query.Where("product_categories.name = ?", params.Name)
	}
	if len(params.IDs) > 0 {
		query = query.Where("product_categories.id IN ?", params.IDs)
	}
	if params.ProductID > 0 {
		query = query.Joins("JOIN product_has_categories phc ON phc.category_id = product_categories.id").
			Where("phc.product_id = ?", params.ProductID)
	}
	return query
}

func (r *categoryRepository) FindAll(params FindAllCategoryDTO) ([]Category, error) {
	var categories []Category

	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit == 0 {
		limit = 10
	}

	query := r.applyFilters(r.db.Model(&Category{}), params)
	orderClause := helpers.ParseOrderBy(params.OrderBy, AllowedOrderColumns, "id DESC")
	query = query.Order("product_categories." + orderClause)

	if limit != -1 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) Count(params FindAllCategoryDTO) (int64, error) {
	var total int64
	query := r.applyFilters(r.db.Model(&Category{}), params)
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *categoryRepository) FindByID(id uint) (*Category, error) {
	var item Category
	if err := r.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *categoryRepository) FindByIDs(ids []uint) ([]Category, error) {
	if len(ids) == 0 {
		return []Category{}, nil
	}
	var items []Category
	if err := r.db.Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *categoryRepository) CountByName(name string, excludeID uint) (int64, error) {
	var count int64
	query := r.db.Model(&Category{}).Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *categoryRepository) CreateRecord(name string, description *string) (*Category, error) {
	item := &Category{
		Name:        name,
		Description: description,
	}
	if err := r.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *categoryRepository) Create(item *Category) error {
	return r.db.Create(item).Error
}

func (r *categoryRepository) Update(item *Category) error {
	return r.db.Save(item).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&Category{}, id).Error
}
