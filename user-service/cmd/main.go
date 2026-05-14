package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/stepup-ai/user-service/internal/cache"
	usergrpc "github.com/stepup-ai/user-service/internal/delivery/grpc"
	"github.com/stepup-ai/user-service/internal/email"
	"github.com/stepup-ai/user-service/internal/messaging"
	"github.com/stepup-ai/user-service/internal/repository"
	"github.com/stepup-ai/user-service/internal/usecase"
	pb "github.com/stepup-ai/user-service/proto/user"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecretKey := os.Getenv("JWT_SECRET")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	grpcPort := os.Getenv("GRPC_PORT")

	if grpcPort == "" {
		grpcPort = "9001"
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	natsPublisher, err := messaging.NewNatsPublisher(natsURL)
	if err != nil {
		log.Fatalf("failed to connect to nats: %v", err)
	}
	defer natsPublisher.Close()

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	redisCache := cache.NewRedisCache(redisURL)
	if err := redisCache.Ping(context.Background()); err != nil {
		log.Printf("warning: redis connection failed: %v", err)
	}

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

	userRepository := repository.NewPostgresUserRepository(db)
	emailSender := email.NewSMTPEmailSender(smtpHost, smtpPort, smtpUser, smtpPass)
	userUsecase := usecase.NewUserUsecase(userRepository, jwtSecretKey, emailSender, natsPublisher, redisCache)
	userHandler := usergrpc.NewUserGRPCHandler(userUsecase)

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", grpcPort, err)
	}

	log.Printf("user service running on port %s", grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
