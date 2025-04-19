package question

import (
	"context"

	pb "github.com/ShopOnGO/review-proto/pkg/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcQuestionService struct {
	pb.UnimplementedQuestionServiceServer
	questionSvc *QuestionService
}

func NewGrpcQuestionService(svc *QuestionService) *GrpcQuestionService {
	return &GrpcQuestionService{questionSvc: svc}
}

func (g *GrpcQuestionService) GetQuestionsForProduct(ctx context.Context, req *pb.GetQuestionsRequest) (*pb.QuestionListResponse, error) {
	questions, err := g.questionSvc.GetQuestionsForProduct(uint(req.ProductVariantId), int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	resp := &pb.QuestionListResponse{}
	for _, q := range questions {
		protoModel := &pb.Model{
			Id:        uint32(q.ID),
			CreatedAt: timestamppb.New(q.CreatedAt),
			UpdatedAt: timestamppb.New(q.UpdatedAt),
			DeletedAt: func() *timestamppb.Timestamp {
				if q.DeletedAt.Valid {
					return timestamppb.New(q.DeletedAt.Time)
				}
				return nil
			}(),
		}

		protoQuestion := &pb.Question{
			Model:            protoModel,
			ProductVariantId: uint32(q.ProductVariantID),
			QuestionText:     q.QuestionText,
			AnswerText:       q.AnswerText,
		}

		if q.UserID != nil {
			protoQuestion.Author = &pb.Question_UserId{UserId: uint32(*q.UserID)}
		} else if len(q.GuestID) > 0 {
			protoQuestion.Author = &pb.Question_GuestId{GuestId: q.GuestID}
		}

		resp.Questions = append(resp.Questions, protoQuestion)
	}

	return resp, nil
}
