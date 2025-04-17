package question

import (
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
)

type QuestionService struct {
	QuestionRepository *QuestionRepository
}

func NewQuestionService(questionRepo *QuestionRepository) *QuestionService {
	return &QuestionService{
		QuestionRepository: questionRepo,
	}
}

func (s *QuestionService) AddQuestion(productVariantID uint, questionText string) (*Question, error) {
	if productVariantID == 0 || questionText == "" {
		return nil, fmt.Errorf("invalid input parameters")
	}

	question := &Question{
		ProductVariantID: productVariantID,
		QuestionText:     questionText,
	}

	if err := s.QuestionRepository.CreateQuestion(question); err != nil {
		logger.Errorf("Error creating question: %v", err)
		return nil, err
	}

	return question, nil
}

func (s *QuestionService) GetQuestionByID(questionID uint) (*Question, error) {
	if questionID == 0 {
		return nil, fmt.Errorf("неверный ID вопроса")
	}

	question, err := s.QuestionRepository.GetQuestionByID(questionID)
	if err != nil {
		logger.Errorf("Ошибка при получении вопроса: %v", err)
		return nil, err
	}
	return question, nil
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


func (s *QuestionService) GetQuestionsForProduct(productVariantID uint, limit, offset int) ([]*Question, error) {
    if productVariantID == 0 {
        return nil, fmt.Errorf("productVariantID is required")
    }

    questions, err := s.QuestionRepository.GetQuestionsByProductVariantIDPaginated(productVariantID, limit, offset)
    if err != nil {
        logger.Errorf("Error getting paginated questions for product %d: %v", productVariantID, err)
        return nil, err
    }

    return questions, nil
}