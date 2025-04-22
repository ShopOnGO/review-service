package review

import (
	"net/http"
	"strconv"
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