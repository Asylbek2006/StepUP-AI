package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/stepup-ai/university-service/internal/entity"
	"github.com/stepup-ai/university-service/internal/repository"
)

const testDatabaseURL = "postgres://postgres:postgres@localhost:5432/university_test_db?sslmode=disable"

func setupTestDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", testDatabaseURL)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	_, err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		t.Fatalf("failed to create extension: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS universities (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			country VARCHAR(100) NOT NULL,
			acceptance_rate DECIMAL(5,2) DEFAULT 0,
			type VARCHAR(50) DEFAULT '',
			category VARCHAR(50) DEFAULT '',
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		t.Fatalf("failed to create universities table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS grants (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			description TEXT DEFAULT '',
			amount DECIMAL(10,2) DEFAULT 0,
			deadline VARCHAR(50) DEFAULT '',
			country VARCHAR(100) DEFAULT '',
			min_gpa DECIMAL(3,2) DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		t.Fatalf("failed to create grants table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS saved_universities (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL,
			university_id UUID REFERENCES universities(id) ON DELETE CASCADE,
			saved_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		t.Fatalf("failed to create saved_universities table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS saved_grants (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL,
			grant_id UUID REFERENCES grants(id) ON DELETE CASCADE,
			saved_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		t.Fatalf("failed to create saved_grants table: %v", err)
	}

	return db
}

func cleanupTestDatabase(db *sql.DB) {
	db.Exec(`DELETE FROM saved_grants`)
	db.Exec(`DELETE FROM saved_universities`)
	db.Exec(`DELETE FROM grants`)
	db.Exec(`DELETE FROM universities`)
}

func TestIntegration_SearchUniversities(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()
	defer cleanupTestDatabase(db)

	_, err := db.Exec(`
		INSERT INTO universities (id, name, country, acceptance_rate, type, category)
		VALUES (uuid_generate_v4(), 'MIT', 'USA', 7.3, 'private', 'reach'),
		       (uuid_generate_v4(), 'Harvard', 'USA', 5.0, 'private', 'reach'),
		       (uuid_generate_v4(), 'Nazarbayev University', 'Kazakhstan', 30.0, 'public', 'target')
	`)
	assert.NoError(t, err)

	universityRepository := repository.NewPostgresUniversityRepository(db)

	universities, err := universityRepository.SearchUniversities(context.Background(), "USA", "", 0, 0)

	assert.NoError(t, err)
	assert.Len(t, universities, 2)
}

func TestIntegration_SaveAndGetUniversity(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()
	defer cleanupTestDatabase(db)

	var universityID string
	err := db.QueryRow(`
		INSERT INTO universities (id, name, country, acceptance_rate, type, category)
		VALUES (uuid_generate_v4(), 'MIT', 'USA', 7.3, 'private', 'reach')
		RETURNING id
	`).Scan(&universityID)
	assert.NoError(t, err)

	universityRepository := repository.NewPostgresUniversityRepository(db)

	savedUniversity := &entity.SavedUniversity{
		ID:           "550e8400-e29b-41d4-a716-446655440000",
		UserID:       "660e8400-e29b-41d4-a716-446655440001",
		UniversityID: universityID,
		SavedAt:      time.Now(),
	}

	err = universityRepository.SaveUniversity(context.Background(), savedUniversity)
	assert.NoError(t, err)

	savedUniversities, err := universityRepository.GetSavedUniversities(context.Background(), "660e8400-e29b-41d4-a716-446655440001")
	assert.NoError(t, err)
	assert.Len(t, savedUniversities, 1)
	assert.Equal(t, "MIT", savedUniversities[0].Name)
}

func TestIntegration_SearchGrants(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()
	defer cleanupTestDatabase(db)

	_, err := db.Exec(`
		INSERT INTO grants (id, name, description, amount, deadline, country, min_gpa)
		VALUES (uuid_generate_v4(), 'Bolashak', 'Government grant', 50000, '2026-12-01', 'Kazakhstan', 3.5),
		       (uuid_generate_v4(), 'Chevening', 'UK grant', 30000, '2026-11-01', 'UK', 3.0)
	`)
	assert.NoError(t, err)

	universityRepository := repository.NewPostgresUniversityRepository(db)

	grants, err := universityRepository.SearchGrants(context.Background(), "Kazakhstan", 3.8)

	assert.NoError(t, err)
	assert.Len(t, grants, 1)
	assert.Equal(t, "Bolashak", grants[0].Name)
}
