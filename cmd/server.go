package main

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/ShopOnGO/review-service/configs"
	"github.com/ShopOnGO/review-service/internal/question"
	"github.com/ShopOnGO/review-service/internal/review"
	"google.golang.org/grpc"

	"github.com/ShopOnGO/review-service/migrations"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/ShopOnGO/review-service/pkg/db"

	pb "github.com/ShopOnGO/review-proto/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func main() {
	migrations.CheckForMigrations()
	conf := configs.LoadConfig()
	database := db.NewDB(conf)

	// repository
	reviewRepo := review.NewReviewRepository(database)
	questionRepo := question.NewQuestionRepository(database)

	// service
	reviewSvc := review.NewReviewService(reviewRepo)
	questionSvc := question.NewQuestionService(questionRepo)

	// Инициализация Kafka-консьюмера
	kafkaConsumer := kafkaService.NewConsumer(
		conf.Kafka.Brokers,
		conf.Kafka.Topic,
		conf.Kafka.GroupID,
		conf.Kafka.ClientID,
	)
	defer kafkaConsumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go kafkaConsumer.Consume(ctx, func(msg kafka.Message) error {
		return handleKafkaMessage(msg, reviewSvc, questionSvc)
	})

	router := gin.Default()

	// handler
	review.NewReviewHandler(router, reviewSvc)
	question.NewQuestionHandler(router, questionSvc)

	go func() {
		if err := router.Run(":8080"); err != nil {
			fmt.Println("Ошибка при запуске HTTP-сервера:", err)
		}
	}()

	go func() {
		listener, err := net.Listen("tcp", ":50052")
		if err != nil {
			fmt.Println("Ошибка при создании TCP listener:", err)
			return
		}

		grpcServer := grpc.NewServer()
		grpcReviewService := review.NewGrpcReviewService(reviewSvc)
		grpcQuestionService := question.NewGrpcQuestionService(questionSvc)

		pb.RegisterReviewServiceServer(grpcServer, grpcReviewService)
		pb.RegisterQuestionServiceServer(grpcServer, grpcQuestionService)

		fmt.Println("gRPC сервер слушает на :50052")
		if err := grpcServer.Serve(listener); err != nil {
			fmt.Println("Ошибка при запуске gRPC сервера:", err)
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
