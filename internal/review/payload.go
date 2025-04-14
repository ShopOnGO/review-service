package review

type BaseReviewEvent struct {
	Action string 				`json:"action"`
	Review ReviewCreatedEvent   `json:"product"`
	UserID uint   				`json:"user_id"`
}

type ReviewCreatedEvent struct {
	ProductVariantID uint   `json:"product_variant_id"`
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