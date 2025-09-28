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

func (s *QuestionService) AddQuestion(productID uint, questionText string, userID *uint, guestID *string) (*Question, error) {
	if productID == 0 || questionText == "" {
		return nil, fmt.Errorf("invalid input parameters")
	}
	var guestIDBytes []byte
	if guestID != nil {
		guestIDBytes = []byte(*guestID)
	}

	question := &Question{
		ProductID:		productID,
		QuestionText:   questionText,
		UserID:         userID,
        GuestID:        guestIDBytes,
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
	_, err := s.QuestionRepository.GetQuestionByID(questionID)
	if err != nil {
		return fmt.Errorf("question with id %d not found", questionID)
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


func (s *QuestionService) GetQuestionsForProduct(productID uint, limit, offset int) ([]*Question, error) {
    if productID == 0 {
        return nil, fmt.Errorf("productID is required")
    }

    questions, err := s.QuestionRepository.GetQuestionsByProductIDPaginated(productID, limit, offset)
    if err != nil {
        logger.Errorf("Error getting paginated questions for product %d: %v", productID, err)
        return nil, err
    }

    return questions, nil
}

func (s *QuestionService) AddLikeToQuestion(questionID, userID uint) (uint, error) {
    if questionID == 0 {
        return 0, fmt.Errorf("invalid question id")
    }

    // Можно здесь добавить дополнительную логику:
    // – проверка, что пользователь ещё не лайкал (если нужна такая фича)

    newCount, err := s.QuestionRepository.IncrementLikes(questionID)
    if err != nil {
        return 0, err
    }
    return newCount, nil
}

func (s *QuestionService) RemoveLikeToQuestion(questionID, userID uint) (uint, error) {
    if questionID == 0 {
        return 0, fmt.Errorf("invalid question id")
    }

    // Можно здесь добавить дополнительную логику:
    // – проверка, что пользователь ещё не лайкал (если нужна такая фича)

    newCount, err := s.QuestionRepository.DecrementLikes(questionID)
    if err != nil {
        return 0, err
    }
    return newCount, nil
}