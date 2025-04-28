package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"

	"github.com/ShopOnGO/review-service/configs"
	"github.com/ShopOnGO/review-service/internal/question"
	"github.com/ShopOnGO/review-service/internal/review"
	"github.com/ShopOnGO/review-service/migrations"
	"github.com/ShopOnGO/review-service/pkg/db"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	pb "github.com/ShopOnGO/review-proto/pkg/service"
)

var (
	httpSrv *http.Server
)

type App struct {
	conf         *configs.Config
	reviewSvc    *review.ReviewService
	questionSvc  *question.QuestionService
	kafkaConsumer *kafkaService.KafkaService
}


func InitServices() *App {
	migrations.CheckForMigrations()
	conf := configs.LoadConfig()
	database := db.NewDB(conf)

	reviewRepo := review.NewReviewRepository(database)
	questionRepo := question.NewQuestionRepository(database)

	reviewSvc := review.NewReviewService(reviewRepo)
	questionSvc := question.NewQuestionService(questionRepo)

	kafkaConsumer := kafkaService.NewConsumer(
		conf.Kafka.Brokers,
		conf.Kafka.Topic,
		conf.Kafka.GroupID,
		conf.Kafka.ClientID,
	)

	return &App{
		conf:          conf,
		reviewSvc:     reviewSvc,
		questionSvc:   questionSvc,
		kafkaConsumer: kafkaConsumer,
	}
}

func RunHTTPServer(app *App) {
	router := gin.Default()
	review.NewReviewHandler(router, app.reviewSvc)
	question.NewQuestionHandler(router, app.questionSvc)

	httpSrv = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	logger.Info("HTTP server listening on :8080")
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Infof("HTTP server error: %v\n", err)
	}
}


func RunGRPCServer(app *App, wg *sync.WaitGroup) *grpc.Server {
	defer wg.Done()
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		logger.Infof("TCP listener error: %v\n", err)
		return nil
	}

	grpcServer := grpc.NewServer()
	pb.RegisterReviewServiceServer(grpcServer, review.NewGrpcReviewService(app.reviewSvc))
	pb.RegisterQuestionServiceServer(grpcServer, question.NewGrpcQuestionService(app.questionSvc))

	logger.Info("gRPC server listening on :50052")
	if err := grpcServer.Serve(listener); err != nil {
		logger.Infof("gRPC server error: %v\n", err)
	}
	return grpcServer
}


func RunKafkaConsumer(ctx context.Context, app *App) {
	defer app.kafkaConsumer.Close()

	dispatcher := kafkaService.NewDispatcher()
	dispatcher.Register("review", func(msg kafka.Message) error {
		return review.HandleReviewEvent(msg.Value, string(msg.Key), app.reviewSvc)
	})
	dispatcher.Register("question", func(msg kafka.Message) error {
		return question.HandleQuestionEvent(msg.Value, string(msg.Key), app.questionSvc)
	})

	logger.Info("Kafka consumer started")
	app.kafkaConsumer.Consume(ctx, dispatcher.Dispatch)
}


func WaitForShutdown(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	logger.Info("Shutdown signal received")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdownCancel()
	if httpSrv != nil {
		logger.Info("Shutting down HTTP server...")
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			logger.Infof("HTTP shutdown error: %v\n", err)
		}
	}
}
