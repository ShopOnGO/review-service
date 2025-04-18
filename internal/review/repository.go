package review

import (
	"github.com/ShopOnGO/review-service/pkg/db"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	Db *db.Db
}

func NewReviewRepository(db *db.Db) *ReviewRepository {
	return &ReviewRepository{
		Db: db,
	}
}

func (r *ReviewRepository) CreateReview(review *Review) error {
	return r.Db.Create(review).Error
}

func (r *ReviewRepository) GetReviewByID(id uint) (*Review, error) {
	var review Review
	err := r.Db.First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) GetReviewsByProductVariantID(productVariantID uint) ([]Review, error) {
	var reviews []Review
	err := r.Db.Where("product_variant_id = ?", productVariantID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewRepository) GetReviewsByProductVariantIDPaginated(productVariantID uint, limit, offset int) ([]*Review, error) {
	var reviews []*Review
	result := r.Db.
		Where("product_variant_id = ?", productVariantID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC"). //сначала новые отзывы, можно изменить
		Find(&reviews)

	if result.Error != nil {
		return nil, result.Error
	}

	return reviews, nil
}

func (r *ReviewRepository) UpdateRating(productVariantID uint, newRating int) error {
    return r.Db.Transaction(func(tx *gorm.DB) error {
        res := tx.Exec(`
            UPDATE product_variants
            SET 
                review_count   = review_count + ?,
                rating_sum     = rating_sum   + ?,
                rating = (rating_sum + ?)::numeric / (review_count + 1)
            WHERE id = ?
        `, 1, newRating, newRating, productVariantID)

        if res.Error != nil {
            return res.Error
        }
        return nil
    })
}

// UpdateRatingDelta — корректируем сумму при update (count не меняется)
func (r *ReviewRepository) UpdateRatingDelta(productVariantID uint, oldRating, newRating int) error {
    delta := newRating - oldRating
    return r.Db.Transaction(func(tx *gorm.DB) error {
        res := tx.Exec(`
            UPDATE product_variants
            SET 
              rating_sum     = rating_sum + ?,
              rating = (rating_sum + ?)::numeric / review_count
            WHERE id = ?`,
            delta, delta, productVariantID,
        )
        return res.Error
    })
}

func (r *ReviewRepository) UpdateRatingDelete(productVariantID uint, oldRating int) error {
    return r.Db.Transaction(func(tx *gorm.DB) error {
        res := tx.Exec(`
            UPDATE product_variants
            SET 
              review_count   = review_count - 1,
              rating_sum     = rating_sum   - ?,
              rating = CASE 
                WHEN review_count > 1 
                  THEN (rating_sum - ?)::numeric / (review_count - 1)
                ELSE 0
              END
            WHERE id = ?`,
            oldRating, oldRating, productVariantID,
        )
        return res.Error
    })
}

func (r *ReviewRepository) UpdateReview(review *Review) error {
	return r.Db.Save(review).Error
}

func (r *ReviewRepository) DeleteReview(review *Review) error {
	return r.Db.Delete(review).Error
}
