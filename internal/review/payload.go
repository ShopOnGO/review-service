package review

type BaseReviewEvent struct {
	Action string `json:"action"`
}

type ReviewCreatedEvent struct {
	Action           string `json:"action"`
	ProductVariantID uint   `json:"product_variant_id"`
	UserID           uint   `json:"user_id"`
	Rating           int16  `json:"rating"`
	Comment          string `json:"comment"`
}

type ReviewUpdatedEvent struct {
	Action   string  `json:"action"`
	ReviewID uint    `json:"review_id"`
	Rating  *int16  `json:"rating,omitempty"`
	Comment *string `json:"comment,omitempty"`
}

type ReviewDeletedEvent struct {
	Action   string `json:"action"`
	ReviewID uint   `json:"review_id"`
}