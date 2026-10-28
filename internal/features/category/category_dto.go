package category

type FindAllCategoryDTO struct {
	Page      int    `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" json:"limit" binding:"omitempty,min=-1,max=100"`
	Search    string `form:"q" json:"q"`
	Name      string `form:"name" json:"name"`
	OrderBy   string `form:"orderBy" json:"orderBy"`
	ProductID uint   `form:"productId" json:"productId"`
	IDs       []uint `form:"ids" json:"ids"`
}

type FindByIDCategoryDTO struct {
	ID uint `uri:"id" json:"id" binding:"required,min=1"`
}

type CreateCategoryDTO struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty"`
}

type UpdateCategoryDTO struct {
	ID          uint    `json:"-"`
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty"`
}
