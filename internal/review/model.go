package review

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	UserID             uint      `gorm:"not null" json:"user_id"`
	ProductVariantID   uint      `gorm:"not null" json:"product_variant_id"`
	Rating             int16     `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment            string	 `gorm:"not null" json:"comment"`
}
