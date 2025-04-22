package main

import (
	"context"
	"sync"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/ShopOnGO/review-service/internal/app"
	"google.golang.org/grpc"
)


func main() {
	services := app.InitServices()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// 1) HTTP
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.RunHTTPServer(services)
	}()

	// 2) gRPC
	var grpcServer *grpc.Server
	wg.Add(1)
	go func() {
		grpcServer = app.RunGRPCServer(services, &wg)
	}()

	// 3) Kafka
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.RunKafkaConsumer(ctx, services)
	}()

	app.WaitForShutdown(cancel)

	if grpcServer != nil {
		logger.Info("Stopping gRPC server…")
		grpcServer.GracefulStop()
	}

	wg.Wait()
	logger.Info("All is stopping")
}
