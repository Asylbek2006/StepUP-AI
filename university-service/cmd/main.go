package main

import (
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

	universitygrpc "github.com/stepup-ai/university-service/internal/delivery/grpc"
	"github.com/stepup-ai/university-service/internal/repository"
	"github.com/stepup-ai/university-service/internal/usecase"
	pb "github.com/stepup-ai/university-service/proto/university"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	grpcPort := os.Getenv("GRPC_PORT")

	if grpcPort == "" {
		grpcPort = "9002"
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
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

	universityRepository := repository.NewPostgresUniversityRepository(db)
	universityUsecase := usecase.NewUniversityUsecase(universityRepository)
	universityHandler := universitygrpc.NewUniversityGRPCHandler(universityUsecase)

	grpcServer := grpc.NewServer()
	pb.RegisterUniversityServiceServer(grpcServer, universityHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", grpcPort, err)
	}

	log.Printf("university service running on port %s", grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
