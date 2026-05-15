# StepUp AI

University admission helper. Microservice-based backend written in Go.

## What it does

Helps students apply to foreign universities — analyzes admission chances, finds matching grants, generates preparation roadmaps and reviews essays using OpenAI.

## Stack

Go 1.25, gRPC, PostgreSQL, Redis, NATS, OpenAI API.

## Services

- **api-gateway** (port 8080) — REST API, JWT auth, routes requests to services via gRPC
- **user-service** (port 9001) — auth, profiles, password reset emails
- **university-service** (port 9002) — university and grant search
- **ai-service** (port 9003) — OpenAI integration for analysis and essay review

NATS is used for events between services (for example, user-service publishes "user.registered" and ai-service listens to it).

## Endpoints

Auth: register, login, logout, refresh, forgot-password, reset-password
Profile: get, update
Universities: search, details, save, list saved
Grants: search, save, list saved
AI: analyze admission, generate roadmap, review essay, match grants, history

Total: 20+ endpoints.

## How to run

You need PostgreSQL, Redis and NATS running locally.

```bash
psql -U postgres -c "CREATE DATABASE user_db;"
psql -U postgres -c "CREATE DATABASE university_db;"
psql -U postgres -c "CREATE DATABASE ai_db;"

nats-server &
redis-server &
```

Then run each service in its own terminal:

```bash
cd user-service && go run cmd/main.go
cd university-service && go run cmd/main.go
cd ai-service && go run cmd/main.go
cd api-gateway && go run cmd/main.go
```

Environment variables you'll need:
- `DATABASE_URL` for each service
- `JWT_SECRET` for user-service and api-gateway
- `OPENAI_API_KEY` for ai-service
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` for user-service email
- `NATS_URL`, `REDIS_URL` where applicable

## Tests

```bash
cd user-service && go test ./...
cd university-service && go test ./...
```

Integration tests in university-service need a running PostgreSQL.

## Team

- Asylbek — api-gateway, user-service
- Ansar — ai-service
- Aray — university-service