package review

import (
	"errors"
	"fmt"

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

func (r *ReviewRepository) GetReviewsByProductID(productID uint) ([]Review, error) {
	var reviews []Review
	err := r.Db.Where("product_id = ?", productID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewRepository) GetReviewsByProductIDPaginated(productID uint, limit, offset int) ([]*Review, error) {
	var reviews []*Review
	result := r.Db.
		Where("product_id = ?", productID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC"). //сначала новые отзывы, можно изменить
		Find(&reviews)

	if result.Error != nil {
		return nil, result.Error
	}

	return reviews, nil
}

func (r *ReviewRepository) getLikesCount(reviewID uint) (uint, error) {
    var likesCount uint
    err := r.Db.Raw(`SELECT likes_count FROM reviews WHERE id = ?`, reviewID).Scan(&likesCount).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return 0, fmt.Errorf("review not found")
        }
        return 0, err
    }

    return likesCount, nil
}

func (r *ReviewRepository) UpdateRating(productID uint, newRating int) error {
    return r.Db.Transaction(func(tx *gorm.DB) error {
        res := tx.Exec(`
            UPDATE products
            SET 
                review_count   = review_count + ?,
                rating_sum     = rating_sum   + ?,
                rating = (rating_sum + ?)::numeric / (review_count + 1)
            WHERE id = ?
        `, 1, newRating, newRating, productID)

        if res.Error != nil {
            return res.Error
        }
        return nil
    })
}

// UpdateRatingDelta — корректируем сумму при update (count не меняется)
func (r *ReviewRepository) UpdateRatingDelta(productID uint, oldRating, newRating int) error {
    delta := newRating - oldRating
    return r.Db.Transaction(func(tx *gorm.DB) error {
        res := tx.Exec(`
            UPDATE products
            SET 
              rating_sum     = rating_sum + ?,
              rating = (rating_sum + ?)::numeric / review_count
            WHERE id = ?`,
            delta, delta, productID,
        )
        return res.Error
    })
}

func (r *ReviewRepository) UpdateRatingDelete(productID uint, oldRating int) error {
    return r.Db.Transaction(func(tx *gorm.DB) error {
        res := tx.Exec(`
            UPDATE products
            SET 
              review_count   = review_count - 1,
              rating_sum     = rating_sum   - ?,
              rating = CASE 
                WHEN review_count > 1 
                  THEN (rating_sum - ?)::numeric / (review_count - 1)
                ELSE 0
              END
            WHERE id = ?`,
            oldRating, oldRating, productID,
        )
        return res.Error
    })
}

func (r *ReviewRepository) IncrementLikes(reviewID uint) (uint, error) {
    var newLikes uint

    err := r.Db.Raw(`
        UPDATE reviews
        SET likes_count = likes_count + 1
        WHERE id = ?
        RETURNING likes_count
    `, reviewID).Scan(&newLikes).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) || newLikes == 0 {
            return 0, fmt.Errorf("review not found")
        }
        return 0, err
    }

    return newLikes, nil
}


func (r *ReviewRepository) DecrementLikes(reviewID uint) (uint, error) {
    currentLikes, err := r.getLikesCount(reviewID)
    if err != nil {
        return 0, err
    }
    if currentLikes < 1 {
        return 0, fmt.Errorf("likes count cannot be less than 1")
    }

    var newLikes uint
    err = r.Db.Raw(`
        UPDATE reviews
        SET likes_count = likes_count - 1
        WHERE id = ?
        RETURNING likes_count
    `, reviewID).Scan(&newLikes).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) || newLikes == 0 {
            return 0, fmt.Errorf("review not found")
        }
        return 0, err
    }

    return newLikes, nil
}

func (r *ReviewRepository) UpdateReview(review *Review) error {
	return r.Db.Save(review).Error
}

func (r *ReviewRepository) DeleteReview(review *Review) error {
	return r.Db.Delete(review).Error
}
