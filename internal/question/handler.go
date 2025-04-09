package question

import (
	"encoding/json"
	"fmt"

	"github.com/ShopOnGO/review-service/pkg/logger"
)

func HandleQuestionEvent(msg []byte, key string, questionSvc *QuestionService) error {
	var base BaseQuestionEvent
	if err := json.Unmarshal(msg, &base); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
	}

	switch base.Action {
	case "created":
		var event QuestionCreatedEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			logger.Errorf("Ошибка десериализации события создания вопроса: %v", err)
			return err
		}
		logger.Infof("Получены данные для создания вопроса: product_variant_id=%d, user_id=%d, question_text=%q",
			event.ProductVariantID, event.UserID, event.QuestionText)
		_, err := questionSvc.AddQuestion(event.ProductVariantID, event.UserID, event.QuestionText)
		if err != nil {
			logger.Errorf("Ошибка при создании вопроса: %v", err)
			return err
		}
		logger.Infof("Вопрос успешно создан для product_variant_id: %d", event.ProductVariantID)

	case "answered":
		var event QuestionAnsweredEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			logger.Errorf("Ошибка десериализации события ответа на вопрос: %v", err)
			return err
		}
		if event.QuestionID == 0 {
			return fmt.Errorf("неверный question_id для ответа")
		}
		if event.AnswerText == "" {
			return fmt.Errorf("answer_text отсутствует")
		}
		if err := questionSvc.AnswerQuestion(event.QuestionID, event.AnswerText); err != nil {
			logger.Errorf("Ошибка при ответе на вопрос: %v", err)
			return err
		}
		logger.Infof("Вопрос успешно отвечен. question_id: %d", event.QuestionID)

	case "deleted":
		var event QuestionDeletedEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			logger.Errorf("Ошибка десериализации события удаления вопроса: %v", err)
			return err
		}
		if event.QuestionID == 0 {
			return fmt.Errorf("неверный question_id для удаления")
		}
		if err := questionSvc.DeleteQuestion(event.QuestionID); err != nil {
			logger.Errorf("Ошибка при удалении вопроса: %v", err)
			return err
		}
		logger.Infof("Вопрос успешно удалён. question_id: %d", event.QuestionID)

	default:
		return fmt.Errorf("неизвестное действие для вопроса: %s", base.Action)
	}
	return nil
}