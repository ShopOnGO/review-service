package review

import (
	"context"
	"github.com/ShopOnGO/review-service/pkg/logger"
	pb "github.com/ShopOnGO/review-proto/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReviewService struct {
	pb.UnimplementedReviewServiceServer
	ReviewRepository *ReviewRepository
}

func NewReviewService(reviewRepo *ReviewRepository) *ReviewService {
	return &ReviewService{
		ReviewRepository: reviewRepo,
	}
}

func (s *ReviewService) AddReview(ctx context.Context, req *pb.AddReviewRequest) (*pb.AddReviewResponse, error) {
	if req.ProductVariantId == 0 || req.UserId == 0 {
		return &pb.AddReviewResponse{
			Success: false,
			Message: "Invalid product_variant_id or user_id",
		}, status.Errorf(codes.InvalidArgument, "Invalid product_variant_id or user_id")
	}

	review := &Review{
		ProductVariantID: uint(req.ProductVariantId),
		UserID:           uint(req.UserId),
		Rating:           int16(req.Rating),
		Comment:          req.Comment,
	}

	err := s.ReviewRepository.CreateReview(review)
	if err != nil {
		logger.Errorf("Error creating review: %v", err)
		return &pb.AddReviewResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "Error creating review: %v", err)
	}

	return &pb.AddReviewResponse{
		Success: true,
		Message: "Review created successfully",
		Review: &pb.Review{
			Model: &pb.Model{
				Id: uint32(review.ID),
			},
			ProductVariantId: uint32(review.ProductVariantID),
			UserId:           uint32(review.UserID),
			Rating:           int32(review.Rating),
			Comment:          review.Comment,
		},
	}, nil
}

func (s *ReviewService) GetReviews(ctx context.Context, req *pb.GetReviewsRequest) (*pb.GetReviewsResponse, error) {
	reviews, err := s.ReviewRepository.GetReviewsByProductVariantID(uint(req.ProductVariantId))
	if err != nil {
		logger.Errorf("Error getting reviews: %v", err)
		return &pb.GetReviewsResponse{
			Reviews: nil,
		}, status.Errorf(codes.Internal, "Error getting reviews: %v", err)
	}

	reviewList := make([]*pb.Review, len(reviews))
	for i, review := range reviews {
		reviewList[i] = &pb.Review{
			Model: &pb.Model{
				Id: uint32(review.ID),
			},
			ProductVariantId: uint32(review.ProductVariantID),
			UserId:           uint32(review.UserID),
			Rating:           int32(review.Rating),
			Comment:          review.Comment,
		}
	}

	return &pb.GetReviewsResponse{
		Reviews: reviewList,
	}, nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*pb.UpdateReviewResponse, error) {
	if req.ReviewId == 0 {
		return &pb.UpdateReviewResponse{
			Success: false,
			Message: "Review ID is required",
		}, status.Errorf(codes.InvalidArgument, "Review ID is required")
	}

	review, err := s.ReviewRepository.GetReviewByID(uint(req.ReviewId))
	if err != nil {
		logger.Errorf("Error getting review: %v", err)
		return &pb.UpdateReviewResponse{
			Success: false,
			Message: "Review not found",
		}, status.Errorf(codes.NotFound, "Review not found")
	}

	if req.Rating != 0 {
		review.Rating = int16(req.Rating)
	}
	if req.Comment != "" {
		review.Comment = req.Comment
	}

	err = s.ReviewRepository.UpdateReview(review)
	if err != nil {
		logger.Errorf("Error updating review: %v", err)
		return &pb.UpdateReviewResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "Error updating review: %v", err)
	}

	return &pb.UpdateReviewResponse{
		Success: true,
		Message: "Review updated successfully",
	}, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewResponse, error) {
	if req.ReviewId == 0 {
		return &pb.DeleteReviewResponse{
			Success: false,
			Message: "Review ID is required",
		}, status.Errorf(codes.InvalidArgument, "Review ID is required")
	}

	review, err := s.ReviewRepository.GetReviewByID(uint(req.ReviewId))
	if err != nil {
		logger.Errorf("Error getting review: %v", err)
		return &pb.DeleteReviewResponse{
			Success: false,
			Message: "Review not found",
		}, status.Errorf(codes.NotFound, "Review not found")
	}

	err = s.ReviewRepository.DeleteReview(review)
	if err != nil {
		logger.Errorf("Error deleting review: %v", err)
		return &pb.DeleteReviewResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "Error deleting review: %v", err)
	}

	return &pb.DeleteReviewResponse{
		Success: true,
		Message: "Review deleted successfully",
	}, nil
}

func (s *ReviewService) GetAverageRating(ctx context.Context, req *pb.GetAverageRatingRequest) (*pb.GetAverageRatingResponse, error) {
	reviews, err := s.ReviewRepository.GetReviewsByProductVariantID(uint(req.ProductVariantId))
	if err != nil {
		logger.Errorf("Error getting reviews: %v", err)
		return nil, status.Errorf(codes.Internal, "Error getting reviews: %v", err)
	}

	var totalRating int16
	for _, review := range reviews {
		totalRating += review.Rating
	}

	averageRating := float64(totalRating) / float64(len(reviews))
	return &pb.GetAverageRatingResponse{
		AverageRating: averageRating,
	}, nil
}
