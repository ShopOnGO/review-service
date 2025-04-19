package question

import (
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	UserID           *uint     `json:"user_id"`
	GuestID   		 []byte    `gorm:"type:bytea;index" json:"guest_id"`
	ProductVariantID uint      `gorm:"not null" json:"product_variant_id"`
	QuestionText     string    `gorm:"not null" json:"question_text"`
	AnswerText       string    `json:"answer_text"`
}
