package review

import (
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
)

type ReviewService struct {
	ReviewRepository *ReviewRepository
}

func NewReviewService(reviewRepo *ReviewRepository) *ReviewService {
	return &ReviewService{
		ReviewRepository: reviewRepo,
	}
}

func (s *ReviewService) AddReview(productVariantID, userID uint, rating int16,likesCount int, comment string) (*Review, error) {
	if productVariantID == 0 || userID == 0 {
		return nil, fmt.Errorf("invalid product_variant_id or user_id")
	}

	review := &Review{
		ProductVariantID: productVariantID,
		UserID:           userID,
		Rating:           rating,
		LikesCount: 	  likesCount,
		Comment:          comment,
	}

	if err := s.ReviewRepository.CreateReview(review); err != nil {
		logger.Errorf("Error creating review: %v", err)
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) GetReviewByID(reviewID uint) (*Review, error) {
	if reviewID == 0 {
		return nil, fmt.Errorf("review ID is required")
	}
	review, err := s.ReviewRepository.GetReviewByID(reviewID)
	if err != nil {
		logger.Errorf("Error getting review by ID: %v", err)
		return nil, fmt.Errorf("review not found")
	}
	return review, nil
}


func (s *ReviewService) UpdateReview(reviewID uint, rating int16, comment string) error {
	if reviewID == 0 {
		return fmt.Errorf("review ID is required")
	}

	review, err := s.ReviewRepository.GetReviewByID(reviewID)
	if err != nil {
		logger.Errorf("Error getting review: %v", err)
		return fmt.Errorf("review not found")
	}

	if rating != 0 {
		review.Rating = rating
	}
	if comment != "" {
		review.Comment = comment
	}

	if err := s.ReviewRepository.UpdateReview(review); err != nil {
		logger.Errorf("Error updating review: %v", err)
		return err
	}

	return nil
}

func (s *ReviewService) DeleteReview(reviewID uint) error {
	if reviewID == 0 {
		return fmt.Errorf("review ID is required")
	}

	review, err := s.ReviewRepository.GetReviewByID(reviewID)
	if err != nil {
		logger.Errorf("Error getting review: %v", err)
		return fmt.Errorf("review not found")
	}

	if err := s.ReviewRepository.DeleteReview(review); err != nil {
		logger.Errorf("Error deleting review: %v", err)
		return err
	}

	return nil
}

func (s *ReviewService) GetReviewsForProduct(productVariantID uint, limit, offset int) ([]*Review, error) {
	if productVariantID == 0 {
		return nil, fmt.Errorf("productVariantID is required")
	}

	reviews, err := s.ReviewRepository.GetReviewsByProductVariantIDPaginated(productVariantID, limit, offset)
	if err != nil {
		logger.Errorf("Error getting paginated reviews: %v", err)
		return nil, err
	}

	return reviews, nil
}

func (s *ReviewService) UpdateRatingAfterCreate(productVariantID uint, rating int16) error {
    return s.ReviewRepository.UpdateRating(productVariantID, int(rating))
}

func (s *ReviewService) UpdateRatingAfterUpdate(productVariantID uint, oldRating, newRating int) error {
    return s.ReviewRepository.UpdateRatingDelta(productVariantID, oldRating, newRating)
}

func (s *ReviewService) UpdateRatingAfterDelete(productVariantID uint, oldRating int) error {
    return s.ReviewRepository.UpdateRatingDelete(productVariantID, oldRating)
}