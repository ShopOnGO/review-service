package main

import (
	"context"
	"net"

	pb "github.com/ShopOnGO/review-proto/pkg/service"
	"github.com/ShopOnGO/review-service/configs"
	"github.com/ShopOnGO/review-service/internal/question"
	"github.com/ShopOnGO/review-service/internal/review"
	"github.com/ShopOnGO/review-service/migrations"
	"github.com/ShopOnGO/review-service/pkg/db"
	"github.com/ShopOnGO/review-service/pkg/kafkaService"
	"github.com/ShopOnGO/review-service/pkg/logger"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

func ReviewApp() *grpc.Server {
	conf := configs.LoadConfig()
	db := db.NewDB(conf)

	// Создаем новый gRPC-сервер
	grpcServer := grpc.NewServer()

	// repositories
	reviewRepo := review.NewReviewRepository(db)
	questionRepo := question.NewQuestionRepository(db)

	// services
	reviewSvc := review.NewReviewService(reviewRepo)
	questionSvc := question.NewQuestionService(questionRepo)

	// registration
	pb.RegisterReviewServiceServer(grpcServer, reviewSvc)
	pb.RegisterQuestionServiceServer(grpcServer, questionSvc)

	return grpcServer
}

func main() {
	migrations.CheckForMigrations()

	// Запускаем gRPC-сервер
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Errorf("Error due conn to tcp: %v", err)
		return
	}
	grpcServer := ReviewApp()

	// --- Kafka Consumer ---
	brokers := []string{"kafka:9092"}
	topic := "review-events"
	groupID := "review-group"

	kafkaConsumer := kafkaService.NewConsumer(brokers, topic, groupID)
	defer kafkaConsumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go kafkaConsumer.Consume(ctx, func(msg kafka.Message) error {
		logger.Infof("📨 Получено сообщение из Kafka: %s", string(msg.Value))

		// 👉 Здесь можно десериализовать и передать в нужный сервис
		// Например: reviewSvc.HandleKafkaMessage(msg.Value)

		return nil
	})


	if err := grpcServer.Serve(listener); err != nil {
		logger.Errorf("Error due starting the gRPC server: %v", err)
	}
	logger.Info("gRPC server is running on :50051")
}
