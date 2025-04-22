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