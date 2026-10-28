package product

import (
	"time"

	"go-gin-boilerplate/internal/features/category"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string              `gorm:"type:varchar(255);index:products_name_index;not null" json:"name"`
	Description *string             `gorm:"type:text" json:"description"`
	Price       int                 `gorm:"not null" json:"price"`
	Cost        *int                `gorm:"type:integer" json:"cost"`
	SKU         *string             `gorm:"type:varchar(100);uniqueIndex" json:"sku"`
	Barcode     *string             `gorm:"type:varchar(100);uniqueIndex" json:"barcode"`
	ImagePath   *string             `gorm:"type:varchar(255)" json:"imagePath"`
	IsActive    bool                `gorm:"default:true;not null" json:"isActive"`
	Stock       int                 `gorm:"default:0;not null" json:"stock"`
	CreatedAt   time.Time           `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time           `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"deletedAt,omitempty"`
	Categories  []category.Category `gorm:"many2many:product_has_categories;joinForeignKey:product_id;joinReferences:category_id" json:"categories,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

// ProductHasCategory merepresentasikan pivot table many-to-many antara Product dan Category.
type ProductHasCategory struct {
	ProductID  uint             `gorm:"primaryKey;column:product_id" json:"productId"`
	CategoryID uint             `gorm:"primaryKey;column:category_id" json:"categoryId"`
	Product    Product          `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Category   category.Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

func (ProductHasCategory) TableName() string {
	return "product_has_categories"
}

var AllowedOrderColumns = map[string]string{
	"id":        "id",
	"name":      "name",
	"price":     "price",
	"stock":     "stock",
	"createdAt": "created_at",
	"updatedAt": "updated_at",
}
