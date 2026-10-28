package category

import (
	"strings"

	services "go-gin-boilerplate/internal/features/category/_services"
	"go-gin-boilerplate/internal/helpers"
)

type CategoryService interface {
	FindAll(params FindAllCategoryDTO) (helpers.PaginateResult[Category], error)
	FindByID(id uint) (*Category, error)
	Create(dto CreateCategoryDTO) (*Category, error)
	Update(dto UpdateCategoryDTO) (*Category, error)
	Delete(id uint) error
}

type categoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) FindAll(params FindAllCategoryDTO) (helpers.PaginateResult[Category], error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	data, err := s.repo.FindAll(params)
	if err != nil {
		return helpers.PaginateResult[Category]{}, err
	}

	total, err := s.repo.Count(params)
	if err != nil {
		return helpers.PaginateResult[Category]{}, err
	}

	return helpers.Paginate(helpers.PaginateParam[Category]{
		Page:  params.Page,
		Limit: params.Limit,
		Total: total,
		Data:  data,
	}), nil
}

func (s *categoryService) FindByID(id uint) (*Category, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, helpers.NewErrorNotFound("Kategori tidak ditemukan")
	}
	return item, nil
}

// Create mendelegasikan logic pembuatan ke _services/category_service_create.go
// namun tetap dipanggil oleh controller melalui CategoryService (category_service.go).
func (s *categoryService) Create(dto CreateCategoryDTO) (*Category, error) {
	return services.CreateCategory[Category](s.repo, dto.Name, dto.Description)
}

func (s *categoryService) Update(dto UpdateCategoryDTO) (*Category, error) {
	existing, err := s.FindByID(dto.ID)
	if err != nil {
		return nil, err
	}

	cleanName := strings.TrimSpace(dto.Name)
	count, err := s.repo.CountByName(cleanName, dto.ID)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, helpers.NewErrorBadRequest("Nama kategori telah digunakan.")
	}

	existing.Name = cleanName
	existing.Description = dto.Description

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *categoryService) Delete(id uint) error {
	if _, err := s.FindByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
