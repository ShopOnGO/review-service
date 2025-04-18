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

	reviewGroup := router.Group("/reviews-service")
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
    logger.Infof("Получено сообщение: %s", string(msg))

    var base BaseReviewEvent
    if err := json.Unmarshal(msg, &base); err != nil {
        return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
    }

    eventHandlers := map[string]func([]byte, *ReviewService) error{
        "create": HandleCreateReviewEvent,
        "update": HandleUpdateReviewEvent,
        "delete": HandleDeleteReviewEvent,
    }

    handler, exists := eventHandlers[base.Action]
    if !exists {
        return fmt.Errorf("неизвестное действие для отзыва: %s", base.Action)
    }

    return handler(msg, reviewSvc)
}


func HandleCreateReviewEvent(msg []byte, reviewSvc *ReviewService) error {
    var base BaseReviewEvent
    if err := json.Unmarshal(msg, &base); err != nil {
        return fmt.Errorf("ошибка десериализации базового сообщения: %w", err)
    }

    event := base.Review
    logger.Infof("Получены данные для создания отзыва: product_variant_id=%d, user_id=%d, rating=%d, comment=%q",
        event.ProductVariantID, base.UserID, event.Rating, event.Comment)

    reviewCreated, err := reviewSvc.AddReview(event.ProductVariantID, base.UserID, event.Rating, event.Comment)
    if err != nil {
        logger.Errorf("Ошибка при создании отзыва: %v", err)
        return err
    }

    if err := reviewSvc.UpdateRatingAfterCreate(reviewCreated.ProductVariantID, reviewCreated.Rating); err != nil {
        logger.Errorf("Ошибка при обновлении агрегатов рейтинга после создания: %v", err)
    }

    logger.Infof("Отзыв успешно создан: %+v", reviewCreated)
    return nil
}

func HandleUpdateReviewEvent(msg []byte, reviewSvc *ReviewService) error {
    var event ReviewUpdatedEvent
    if err := json.Unmarshal(msg, &event); err != nil {
        logger.Errorf("Ошибка десериализации события обновления отзыва: %v", err)
        return err
    }

    oldReview, err := reviewSvc.GetReviewByID(event.ReviewID)
    if err != nil {
        return err
    }

    var newRating int16 = oldReview.Rating
    if event.Rating != nil {
        newRating = *event.Rating
    }
    var newComment string = oldReview.Comment
    if event.Comment != nil {
        newComment = *event.Comment
    }

    if err := reviewSvc.UpdateReview(event.ReviewID, newRating, newComment); err != nil {
        logger.Errorf("Ошибка при обновлении отзыва: %v", err)
        return err
    }

    if event.Rating != nil {
        if err := reviewSvc.UpdateRatingAfterUpdate(oldReview.ProductVariantID, int(oldReview.Rating), int(newRating)); err != nil {
            logger.Errorf("Ошибка при обновлении агрегатов рейтинга после редактирования: %v", err)
        }
    }

    logger.Infof("Отзыв успешно обновлён. review_id: %d", event.ReviewID)
    return nil
}

func HandleDeleteReviewEvent(msg []byte, reviewSvc *ReviewService) error {
    var event ReviewDeletedEvent
    if err := json.Unmarshal(msg, &event); err != nil {
        logger.Errorf("Ошибка десериализации события удаления отзыва: %v", err)
        return err
    }

    oldReview, err := reviewSvc.GetReviewByID(event.ReviewID)
    if err != nil {
        return err
    }

    if err := reviewSvc.DeleteReview(event.ReviewID); err != nil {
        logger.Errorf("Ошибка при удалении отзыва: %v", err)
        return err
    }

    if err := reviewSvc.UpdateRatingAfterDelete(oldReview.ProductVariantID, int(oldReview.Rating)); err != nil {
        logger.Errorf("Ошибка при обновлении агрегатов рейтинга после удаления: %v", err)
    }

    logger.Infof("Отзыв успешно удалён. review_id: %d", event.ReviewID)
    return nil
}