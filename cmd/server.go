package main

import (
	"net"

	pb "github.com/ShopOnGO/review-proto/pkg/service"
	"github.com/ShopOnGO/review-service/configs"
	"github.com/ShopOnGO/review-service/migrations"
	"github.com/ShopOnGO/review-service/pkg/db"
	"github.com/ShopOnGO/review-service/pkg/logger"

	"github.com/ShopOnGO/review-service/internal/review"
	"github.com/ShopOnGO/review-service/internal/question"

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
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Errorf("Error due conn to tcp: %v", err)
		return
	}
	migrations.CheckForMigrations()
	logger.Info("gRPC server is running on :50051")

	grpcServer := ReviewApp()

	if err := grpcServer.Serve(listener); err != nil {
		logger.Errorf("Error due starting the gRPC server: %v", err)
	}
}
