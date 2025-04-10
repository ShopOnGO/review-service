package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ShopOnGO/review-service/configs"
	"github.com/ShopOnGO/review-service/internal/question"
	"github.com/ShopOnGO/review-service/internal/review"
	"github.com/ShopOnGO/review-service/migrations"
	"github.com/ShopOnGO/review-service/pkg/db"
	"github.com/ShopOnGO/review-service/pkg/kafkaService"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func main() {
	migrations.CheckForMigrations()
	conf := configs.LoadConfig()
	database := db.NewDB(conf)

	// репозитории
	reviewRepo := review.NewReviewRepository(database)
	questionRepo := question.NewQuestionRepository(database)

	// сервисы
	reviewSvc := review.NewReviewService(reviewRepo)
	questionSvc := question.NewQuestionService(questionRepo)

	// Инициализация Kafka-консьюмера
	kafkaConsumer := kafkaService.NewConsumer(
		conf.Kafka.Brokers,
		conf.Kafka.Topic,
		conf.Kafka.GroupID,
	)
	defer kafkaConsumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go kafkaConsumer.Consume(ctx, func(msg kafka.Message) error {
		return handleKafkaMessage(msg, reviewSvc, questionSvc)
	})

	router := gin.Default()
	review.NewReviewHandler(router, reviewSvc)
	question.NewQuestionHandler(router, questionSvc)

	go func() {
		if err := router.Run(":8080"); err != nil {
			fmt.Println("Ошибка при запуске HTTP-сервера:", err)
		}
	}()

	select {}
}

func handleKafkaMessage(msg kafka.Message, reviewSvc *review.ReviewService, questionSvc *question.QuestionService) error {
	key := string(msg.Key)

	switch {
	case strings.HasPrefix(key, "review-"):
		return review.HandleReviewEvent(msg.Value, key, reviewSvc)
	case strings.HasPrefix(key, "question-"):
		return question.HandleQuestionEvent(msg.Value, key, questionSvc)
	default:
		return fmt.Errorf("неподдерживаемый ключ: %s", key)
	}
}