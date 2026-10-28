package product

type FindAllProductDTO struct {
	Page       int    `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit      int    `form:"limit" json:"limit" binding:"omitempty,min=-1,max=100"`
	Search     string `form:"q" json:"q"`
	SKU        string `form:"sku" json:"sku"`
	Barcode    string `form:"barcode" json:"barcode"`
	CategoryID uint   `form:"categoryId" json:"categoryId"`
	IsActive   *bool  `form:"isActive" json:"isActive"`
	OrderBy    string `form:"orderBy" json:"orderBy"`
	IDs        []uint `form:"ids" json:"ids"`
}

type FindByIDProductDTO struct {
	ID uint `uri:"id" json:"id" binding:"required,min=1"`
}

type CreateProductDTO struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty"`
	Price       int     `json:"price" binding:"required,min=0"`
	Cost        *int    `json:"cost" binding:"omitempty,min=0"`
	SKU         *string `json:"sku" binding:"omitempty,max=100"`
	Barcode     *string `json:"barcode" binding:"omitempty,max=100"`
	ImagePath   *string `json:"imagePath" binding:"omitempty"`
	IsActive    *bool   `json:"isActive" binding:"omitempty"`
	Stock       int     `json:"stock" binding:"omitempty,min=0"`
	CategoryIDs []uint  `json:"categoryIds" binding:"omitempty"`
}

type UpdateProductDTO struct {
	ID          uint    `json:"-"`
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty"`
	Price       int     `json:"price" binding:"required,min=0"`
	Cost        *int    `json:"cost" binding:"omitempty,min=0"`
	SKU         *string `json:"sku" binding:"omitempty,max=100"`
	Barcode     *string `json:"barcode" binding:"omitempty,max=100"`
	ImagePath   *string `json:"imagePath" binding:"omitempty"`
	IsActive    *bool   `json:"isActive" binding:"omitempty"`
	Stock       int     `json:"stock" binding:"omitempty,min=0"`
	CategoryIDs *[]uint `json:"categoryIds" binding:"omitempty"`
}
