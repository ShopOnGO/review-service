package question

import (
	"context"

	pb "github.com/ShopOnGO/review-proto/pkg/service"
	"github.com/ShopOnGO/review-service/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type QuestionService struct {
	pb.UnimplementedQuestionServiceServer
	QuestionRepository *QuestionRepository
}

func NewQuestionService(questionRepo *QuestionRepository) *QuestionService {
	return &QuestionService{
		QuestionRepository: questionRepo,
	}
}

func (s *QuestionService) AddQuestion(ctx context.Context, req *pb.AddQuestionRequest) (*pb.AddQuestionResponse, error) {
	if req.ProductVariantId == 0 || req.UserId == 0 || req.QuestionText == "" {
		return &pb.AddQuestionResponse{
			Success: false,
			Message: "Invalid input parameters",
		}, status.Errorf(codes.InvalidArgument, "Invalid input parameters")
	}

	question := &Question{
		ProductVariantID: uint(req.ProductVariantId),
		UserID:           uint(req.UserId),
		QuestionText:     req.QuestionText,
	}

	err := s.QuestionRepository.CreateQuestion(question)
	if err != nil {
		logger.Errorf("Error creating question: %v", err)
		return &pb.AddQuestionResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "Error creating question: %v", err)
	}

	return &pb.AddQuestionResponse{
		Success: true,
		Message: "Question added successfully",
		Question: &pb.Question{
			Model: &pb.Model{
				Id: uint32(question.ID),
			},
			ProductVariantId: uint32(question.ProductVariantID),
			UserId:           uint32(question.UserID),
			QuestionText:     question.QuestionText,
			AnswerText:       question.AnswerText,
		},
	}, nil
}

func (s *QuestionService) GetQuestions(ctx context.Context, req *pb.GetQuestionsRequest) (*pb.GetQuestionsResponse, error) {
	questions, err := s.QuestionRepository.GetQuestionsByProductVariantID(uint(req.ProductVariantId))
	if err != nil {
		logger.Errorf("Error getting questions: %v", err)
		return &pb.GetQuestionsResponse{
			Questions: nil,
		}, status.Errorf(codes.Internal, "Error getting questions: %v", err)
	}

	questionList := make([]*pb.Question, len(questions))
	for i, question := range questions {
		questionList[i] = &pb.Question{
			Model: &pb.Model{
				Id: uint32(question.ID),
			},
			ProductVariantId: uint32(question.ProductVariantID),
			UserId:           uint32(question.UserID),
			QuestionText:     question.QuestionText,
			AnswerText:       question.AnswerText,
		}
	}

	return &pb.GetQuestionsResponse{
		Questions: questionList,
	}, nil
}


func (s *QuestionService) AnswerQuestion(ctx context.Context, req *pb.AnswerQuestionRequest) (*pb.AnswerQuestionResponse, error) {
	if req.QuestionId == 0 || req.AnswerText == "" {
		return &pb.AnswerQuestionResponse{
			Success: false,
			Message: "Invalid input parameters",
		}, status.Errorf(codes.InvalidArgument, "Invalid input parameters")
	}

	err := s.QuestionRepository.UpdateAnswer(uint(req.QuestionId), req.AnswerText)
	if err != nil {
		logger.Errorf("Error answering question: %v", err)
		return &pb.AnswerQuestionResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "Error answering question: %v", err)
	}

	return &pb.AnswerQuestionResponse{
		Success: true,
		Message: "Answer added successfully",
	}, nil
}


func (s *QuestionService) DeleteQuestion(ctx context.Context, req *pb.DeleteQuestionRequest) (*pb.DeleteQuestionResponse, error) {
	if req.QuestionId == 0 {
		return &pb.DeleteQuestionResponse{
			Success: false,
			Message: "Invalid question ID",
		}, status.Errorf(codes.InvalidArgument, "Invalid question ID")
	}

	err := s.QuestionRepository.DeleteQuestionByID(uint(req.QuestionId))
	if err != nil {
		logger.Errorf("Error deleting question: %v", err)
		return &pb.DeleteQuestionResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "Error deleting question: %v", err)
	}

	return &pb.DeleteQuestionResponse{
		Success: true,
		Message: "Question deleted successfully",
	}, nil
}
