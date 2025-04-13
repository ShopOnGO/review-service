package review

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewSvc *ReviewService
}

func NewReviewHandler(router *gin.Engine,reviewSvc *ReviewService) *ReviewHandler {
	handler := &ReviewHandler{reviewSvc: reviewSvc}

	reviewGroup := router.Group("/review-service")
	{
		reviewGroup.GET("/:id", handler.getReviewByID)
	}

	return handler
}

func (h *ReviewHandler) getReviewByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	review, err := h.reviewSvc.GetReviewByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отзыв не найден"})
		return
	}

	c.JSON(http.StatusOK, review)
}


func HandleReviewEvent(msg []byte, key string, reviewSvc *ReviewService) error {
	
	var base BaseReviewEvent
	if err := json.Unmarshal(msg, &base); err != nil {
		return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
	}

	switch base.Action {
	case "created":
		var event ReviewCreatedEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			logger.Errorf("Ошибка десериализации события создания отзыва: %v", err)
			return err
		}
		logger.Infof("Получены данные для создания отзыва: product_variant_id=%d, user_id=%d, rating=%d, comment=%q",
			event.ProductVariantID, event.UserID, event.Rating, event.Comment)
		reviewCreated, err := reviewSvc.AddReview(event.ProductVariantID, event.UserID, event.Rating, event.Comment)
		if err != nil {
			logger.Errorf("Ошибка при создании отзыва: %v", err)
			return err
		}
		logger.Infof("Отзыв успешно создан: %+v", reviewCreated)
	case "updated":
		var event ReviewUpdatedEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			logger.Errorf("Ошибка десериализации события обновления отзыва: %v", err)
			return err
		}
		var rating int16
		if event.Rating != nil {
			rating = *event.Rating
		}
		var comment string
		if event.Comment != nil {
			comment = *event.Comment
		}
		if err := reviewSvc.UpdateReview(event.ReviewID, rating, comment); err != nil {
			logger.Errorf("Ошибка при обновлении отзыва: %v", err)
			return err
		}
		logger.Infof("Отзыв успешно обновлён. review_id: %d", event.ReviewID)
	case "deleted":
		var event ReviewDeletedEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			logger.Errorf("Ошибка десериализации события удаления отзыва: %v", err)
			return err
		}
		if err := reviewSvc.DeleteReview(event.ReviewID); err != nil {
			logger.Errorf("Ошибка при удалении отзыва: %v", err)
			return err
		}
		logger.Infof("Отзыв успешно удалён. review_id: %d", event.ReviewID)
	default:
		return fmt.Errorf("неизвестное действие для отзыва: %s", base.Action)
	}
	return nil
}
