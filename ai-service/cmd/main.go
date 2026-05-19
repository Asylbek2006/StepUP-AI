package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	aigrpc "github.com/stepup-ai/ai-service/internal/delivery/grpc"
	"github.com/stepup-ai/ai-service/internal/messaging"
	"github.com/stepup-ai/ai-service/internal/repository"
	"github.com/stepup-ai/ai-service/internal/usecase"
	pb "github.com/stepup-ai/ai-service/proto/ai"
	gogrpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	grpcPort := os.Getenv("GRPC_PORT")

	if grpcPort == "" {
		grpcPort = "9003"
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	natsSubscriber, err := messaging.NewNatsSubscriber(natsURL)
	if err != nil {
		log.Printf("warning: nats connection failed: %v", err)
	} else {
		defer natsSubscriber.Close()

		err = natsSubscriber.SubscribeToUserRegisteredEvents(func(event messaging.UserRegisteredEvent) error {
			log.Printf("new user registered: %s (%s)", event.FullName, event.Email)
			return nil
		})
		if err != nil {
			log.Printf("warning: failed to subscribe to user registered events: %v", err)
		}
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	go func() {
		metricsHttpMux := http.NewServeMux()
		metricsHttpMux.Handle("/metrics", promhttp.Handler())
		log.Printf("metrics server running on port 9103")
		if err := http.ListenAndServe(":9103", metricsHttpMux); err != nil {
			log.Printf("metrics server failed: %v", err)
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create migration driver: %v", err)
	}

	migrationRunner, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("failed to create migration runner: %v", err)
	}

	if err := migrationRunner.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %v", err)
	}

	aiRepository := repository.NewPostgresAIRepository(db)
	aiUsecase := usecase.NewAIUsecase(aiRepository, openaiAPIKey)
	aiHandler := aigrpc.NewAIGRPCHandler(aiUsecase)

	grpcServer := grpc.NewServer()
	pb.RegisterAIServiceServer(grpcServer, aiHandler)
	reflection.Register(grpcServer)
	gogrpc_health_v1.RegisterHealthServer(grpcServer, aigrpc.NewHealthCheckHandler())

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", grpcPort, err)
	}

	log.Printf("=== AI Service ===")
	log.Printf("Version: 1.0.0")
	log.Printf("gRPC Port: %s", grpcPort)
	log.Printf("Database: connected")
	log.Printf("NATS: connected")
	log.Printf("==================")

	log.Printf("ai service running on port %s", grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
