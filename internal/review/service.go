package review

import (
	"fmt"
	"github.com/ShopOnGO/review-service/pkg/logger"
)

type ReviewService struct {
	ReviewRepository *ReviewRepository
}

func NewReviewService(reviewRepo *ReviewRepository) *ReviewService {
	return &ReviewService{
		ReviewRepository: reviewRepo,
	}
}

func (s *ReviewService) AddReview(productVariantID, userID uint, rating int16, comment string) (*Review, error) {
	if productVariantID == 0 || userID == 0 {
		return nil, fmt.Errorf("invalid product_variant_id or user_id")
	}

	review := &Review{
		ProductVariantID: productVariantID,
		UserID:           userID,
		Rating:           rating,
		Comment:          comment,
	}

	if err := s.ReviewRepository.CreateReview(review); err != nil {
		logger.Errorf("Error creating review: %v", err)
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) GetReviews(productVariantID uint) ([]Review, error) {
	reviews, err := s.ReviewRepository.GetReviewsByProductVariantID(productVariantID)
	if err != nil {
		logger.Errorf("Error getting reviews: %v", err)
		return nil, err
	}
	return reviews, nil
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

func (s *ReviewService) GetAverageRating(productVariantID uint) (float64, error) {
	reviews, err := s.ReviewRepository.GetReviewsByProductVariantID(productVariantID)
	if err != nil {
		logger.Errorf("Error getting reviews: %v", err)
		return 0, err
	}

	if len(reviews) == 0 {
		return 0, nil
	}

	var total int
	for _, r := range reviews {
		total += int(r.Rating)
	}

	average := float64(total) / float64(len(reviews))
	return average, nil
}
