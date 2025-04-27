package question

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionSvc *QuestionService
}

func NewQuestionHandler(router *gin.Engine, questionSvc *QuestionService) *QuestionHandler {
	handler := &QuestionHandler{questionSvc: questionSvc}

	questionGroup := router.Group("/reviews-service/questions")
	{
		questionGroup.GET("/:id", handler.GetQuestionByID)
	}

	return handler
}

// GetQuestionByID godoc
// @Summary Получить вопрос по ID
// @Description Возвращает вопрос по его уникальному идентификатору
// @Tags Вопросы
// @Param id path int true "ID вопроса"
// @Success 200 {object} question.Question
// @Failure 400 {object} gin.H "Некорректный ID"
// @Failure 404 {object} gin.H "Вопрос не найден"
// @Router /reviews-service/questions/{id} [get]
func (h *QuestionHandler) GetQuestionByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	question, err := h.questionSvc.GetQuestionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Вопрос не найден"})
		return
	}

	c.JSON(http.StatusOK, question)
}