package question

import (
	"fmt"

	"github.com/ShopOnGO/review-service/pkg/logger"
)

type QuestionService struct {
	QuestionRepository *QuestionRepository
}

func NewQuestionService(questionRepo *QuestionRepository) *QuestionService {
	return &QuestionService{
		QuestionRepository: questionRepo,
	}
}

func (s *QuestionService) AddQuestion(productVariantID, userID uint, questionText string) (*Question, error) {
	if productVariantID == 0 || userID == 0 || questionText == "" {
		return nil, fmt.Errorf("invalid input parameters")
	}

	question := &Question{
		ProductVariantID: productVariantID,
		UserID:           userID,
		QuestionText:     questionText,
	}

	if err := s.QuestionRepository.CreateQuestion(question); err != nil {
		logger.Errorf("Error creating question: %v", err)
		return nil, err
	}

	return question, nil
}

func (s *QuestionService) GetQuestions(productVariantID uint) ([]Question, error) {
	questions, err := s.QuestionRepository.GetQuestionsByProductVariantID(productVariantID)
	if err != nil {
		logger.Errorf("Error getting questions: %v", err)
		return nil, err
	}
	return questions, nil
}

func (s *QuestionService) AnswerQuestion(questionID uint, answerText string) error {
	if questionID == 0 || answerText == "" {
		return fmt.Errorf("invalid input parameters")
	}
	if err := s.QuestionRepository.UpdateAnswer(questionID, answerText); err != nil {
		logger.Errorf("Error answering question: %v", err)
		return err
	}
	return nil
}

func (s *QuestionService) DeleteQuestion(questionID uint) error {
	if questionID == 0 {
		return fmt.Errorf("invalid question ID")
	}
	if err := s.QuestionRepository.DeleteQuestionByID(questionID); err != nil {
		logger.Errorf("Error deleting question: %v", err)
		return err
	}
	return nil
}
