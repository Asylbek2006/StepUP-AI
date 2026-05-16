.PHONY: help build test run-user run-university run-ai run-gateway clean

help:
  @echo "Available commands:"
  @echo "  make build           - Build all services"
  @echo "  make test            - Run all tests"
  @echo "  make run-user        - Run user-service"
  @echo "  make run-university  - Run university-service"
  @echo "  make run-ai          - Run ai-service"
  @echo "  make run-gateway     - Run api-gateway"
  @echo "  make clean           - Clean build artifacts"

build:
  cd user-service && go build -o bin/user-service ./cmd
  cd university-service && go build -o bin/university-service ./cmd
  cd ai-service && go build -o bin/ai-service ./cmd
  cd api-gateway && go build -o bin/api-gateway ./cmd

test:
  cd user-service && go test ./...
  cd university-service && go test ./...
  cd ai-service && go test ./...
  cd api-gateway && go test ./...

run-user:
  cd user-service && go run cmd/main.go

run-university:
  cd university-service && go run cmd/main.go

run-ai:
  cd ai-service && go run cmd/main.go

run-gateway:
  cd api-gateway && go run cmd/main.go

clean:
  rm -rf user-service/bin
  rm -rf university-service/bin
  rm -rf ai-service/bin
  rm -rf api-gateway/bin