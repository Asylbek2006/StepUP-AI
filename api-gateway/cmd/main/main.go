package main

import (
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/stepup-ai/api-gateway/internal/handler"
	"github.com/stepup-ai/api-gateway/internal/router"
)

func main() {
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	universityServiceAddr := os.Getenv("UNIVERSITY_SERVICE_ADDR")
	aiServiceAddr := os.Getenv("AI_SERVICE_ADDR")
	jwtSecretKey := os.Getenv("JWT_SECRET")
	httpPort := os.Getenv("HTTP_PORT")

	if httpPort == "" {
		httpPort = "8080"
	}
	if userServiceAddr == "" {
		userServiceAddr = "localhost:9001"
	}
	if universityServiceAddr == "" {
		universityServiceAddr = "localhost:9002"
	}
	if aiServiceAddr == "" {
		aiServiceAddr = "localhost:9003"
	}

	userServiceConnection, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	defer userServiceConnection.Close()

	universityServiceConnection, err := grpc.NewClient(universityServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to university service: %v", err)
	}
	defer universityServiceConnection.Close()

	aiServiceConnection, err := grpc.NewClient(aiServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to ai service: %v", err)
	}
	defer aiServiceConnection.Close()

	userHandler := handler.NewUserHandler(userServiceConnection)
	universityHandler := handler.NewUniversityHandler(universityServiceConnection)
	aiHandler := handler.NewAIHandler(aiServiceConnection)

	appRouter := router.SetupRouter(userHandler, universityHandler, aiHandler, jwtSecretKey)

	log.Printf("api gateway running on port %s", httpPort)
	if err := appRouter.Run(":" + httpPort); err != nil {
		log.Fatalf("failed to run api gateway: %v", err)
	}
}
