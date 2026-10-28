package category

import "time"

type Category struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Category) TableName() string {
	return "product_categories"
}

var AllowedOrderColumns = map[string]string{
	"id":        "id",
	"name":      "name",
	"createdAt": "created_at",
	"updatedAt": "updated_at",
}
