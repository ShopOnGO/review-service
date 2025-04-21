package review

import (
	"context"

	pb "github.com/ShopOnGO/review-proto/pkg/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcReviewService struct {
	pb.UnimplementedReviewServiceServer
	reviewSvc *ReviewService
}

func NewGrpcReviewService(svc *ReviewService) *GrpcReviewService {
	return &GrpcReviewService{reviewSvc: svc}
}

func (g *GrpcReviewService) GetReviewsForProduct(ctx context.Context, req *pb.GetReviewsRequest) (*pb.ReviewListResponse, error) {
	reviews, err := g.reviewSvc.GetReviewsForProduct(uint(req.ProductVariantId), int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	resp := &pb.ReviewListResponse{}
	for _, r := range reviews {
		resp.Reviews = append(resp.Reviews, &pb.Review{
			Model: &pb.Model{
			  Id:        uint32(r.ID),
			  CreatedAt: timestamppb.New(r.CreatedAt),
			  UpdatedAt: timestamppb.New(r.UpdatedAt),
			  DeletedAt: func() *timestamppb.Timestamp {
				if r.DeletedAt.Valid {
					return timestamppb.New(r.DeletedAt.Time)
				}
				return nil
			}(),
			},
			
			ProductVariantId: uint32(r.ProductVariantID),
			UserId:           uint32(r.UserID),
			Rating:           int32(r.Rating),
			LikesCount: 	  int32(r.LikesCount),
			Comment:          r.Comment,
		  })
		  
	}
	return resp, nil
}
