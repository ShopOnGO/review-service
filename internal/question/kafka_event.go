package question

import (
	"encoding/json"
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
)

func HandleQuestionEvent(msg []byte, key string, questionSvc *QuestionService) error {
	var base BaseQuestionEvent
	if err := json.Unmarshal(msg, &base); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
	}

	eventHandlers := map[string]func([]byte, *QuestionService) error{
		"create":  		HandleCreateQuestionEvent,
		"answer":  		HandleAnswerQuestionEvent,
		"delete":  		HandleDeleteQuestionEvent,
		"addLike": 		HandleAddLikeQuestionEvent,
		"removeLike":	HandleRemoveLikeQuestionEvent,
	}

	handler, exists := eventHandlers[base.Action]
	if !exists {
		return fmt.Errorf("неизвестное действие для вопроса: %s", base.Action)
	}

	return handler(msg, questionSvc)
}

func HandleCreateQuestionEvent(msg []byte, questionSvc *QuestionService) error {
	var event QuestionCreatedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		logger.Errorf("Ошибка десериализации события создания вопроса: %v", err)
		return err
	}
	logger.Infof("Создаём вопрос: variant=%d, text=%q, user=%v, guest=%v",
		event.ProductVariantID, event.QuestionText, event.Author.UserID, event.Author.GuestID)

	_, err := questionSvc.AddQuestion(event.ProductVariantID, event.QuestionText, event.Author.UserID, event.Author.GuestID)
	if err != nil {
		logger.Errorf("Ошибка при создании вопроса: %v", err)
		return err
	}

	logger.Infof("Вопрос успешно создан для product_variant_id: %d", event.ProductVariantID)
	return nil
}

func HandleAnswerQuestionEvent(msg []byte, questionSvc *QuestionService) error {
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
	return nil
}

func HandleDeleteQuestionEvent(msg []byte, questionSvc *QuestionService) error {
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
	return nil
}

func HandleAddLikeQuestionEvent(msg []byte, questionSvc *QuestionService) error {
	logger.Infof("Получено сообщение для лайка: %s", string(msg))

	var event struct {
		QuestionID uint   `json:"question_id"`
		UserID     uint   `json:"user_id"`
	}
	
	if err := json.Unmarshal(msg, &event); err != nil {
		logger.Errorf("Ошибка десериализации события лайка вопроса: %v", err)
		return err
	}

	if event.QuestionID == 0 {
		return fmt.Errorf("неверный question_id для лайка")
	}
	if event.UserID == 0 {
		return fmt.Errorf("неверный user_id для лайка")
	}

	newLikes, err := questionSvc.AddLikeToQuestion(event.QuestionID, event.UserID)
	if err != nil {
		logger.Errorf("Ошибка при добавлении лайка к вопросу: %v", err)
		return err
	}

	logger.Infof("Лайк успешно добавлен. question_id: %d, user_id: %d, new_likes: %d", event.QuestionID, event.UserID, newLikes)
	return nil
}

func HandleRemoveLikeQuestionEvent(msg []byte, questionSvc *QuestionService) error {
	logger.Infof("Получено сообщение для удаления лайка: %s", string(msg))

	var event struct {
		QuestionID uint   `json:"question_id"`
		UserID     uint   `json:"user_id"`
	}
	
	if err := json.Unmarshal(msg, &event); err != nil {
		logger.Errorf("Ошибка десериализации события удаления лайка на вопроса: %v", err)
		return err
	}
	logger.Infof("Удаляем лайк у вопроса: review_id=%d, от user_id=%d", event.QuestionID, event.UserID)

	if event.QuestionID == 0 {
		return fmt.Errorf("неверный question_id для лайка")
	}
	if event.UserID == 0 {
		return fmt.Errorf("неверный user_id для лайка")
	}

	newLikes, err := questionSvc.RemoveLikeToQuestion(event.QuestionID, event.UserID)
	if err != nil {
		logger.Errorf("Ошибка при удалении лайка к вопросу: %v", err)
		return err
	}

	logger.Infof("Лайк успешно удален. question_id: %d, user_id: %d, new_likes: %d", event.QuestionID, event.UserID, newLikes)
	return nil
}
