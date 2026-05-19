package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/stepup-ai/university-service/internal/entity"
)

type UniversityRepository interface {
	SearchUniversities(ctx context.Context, country, universityType string, minAcceptanceRate, maxAcceptanceRate float32) ([]*entity.University, error)
	GetUniversityByID(ctx context.Context, universityID string) (*entity.University, error)
	SaveUniversity(ctx context.Context, savedUniversity *entity.SavedUniversity) error
	GetSavedUniversities(ctx context.Context, userID string) ([]*entity.University, error)
	SearchGrants(ctx context.Context, country string, gpa float32) ([]*entity.Grant, error)
	GetGrantByID(ctx context.Context, grantID string) (*entity.Grant, error)
	SaveGrant(ctx context.Context, savedGrant *entity.SavedGrant) error
	GetSavedGrants(ctx context.Context, userID string) ([]*entity.Grant, error)
	RemoveSavedUniversity(ctx context.Context, userID, universityID string) error
	RemoveSavedGrant(ctx context.Context, userID, grantID string) error
	ListUniversitiesByCategory(ctx context.Context, category string) ([]*entity.University, error)
	GetUniversityStatistics(ctx context.Context) (int32, int32, int32, int32, error)
}

type PostgresUniversityRepository struct {
	db *sql.DB
}

func NewPostgresUniversityRepository(db *sql.DB) UniversityRepository {
	return &PostgresUniversityRepository{db: db}
}

func (r *PostgresUniversityRepository) SearchUniversities(ctx context.Context, country, universityType string, minAcceptanceRate, maxAcceptanceRate float32) ([]*entity.University, error) {
	query := `
		SELECT id, name, country, acceptance_rate, type, category, created_at
		FROM universities
		WHERE ($1 = '' OR country = $1)
		AND ($2 = '' OR type = $2)
		AND acceptance_rate >= $3
		AND ($4 = 0 OR acceptance_rate <= $4)
	`
	rows, err := r.db.QueryContext(ctx, query, country, universityType, minAcceptanceRate, maxAcceptanceRate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var universities []*entity.University
	for rows.Next() {
		university := &entity.University{}
		if err := rows.Scan(
			&university.ID, &university.Name, &university.Country,
			&university.AcceptanceRate, &university.Type, &university.Category,
			&university.CreatedAt,
		); err != nil {
			return nil, err
		}
		universities = append(universities, university)
	}
	return universities, nil
}

func (r *PostgresUniversityRepository) GetUniversityByID(ctx context.Context, universityID string) (*entity.University, error) {
	query := `SELECT id, name, country, acceptance_rate, type, category, created_at FROM universities WHERE id = $1`
	university := &entity.University{}
	err := r.db.QueryRowContext(ctx, query, universityID).Scan(
		&university.ID, &university.Name, &university.Country,
		&university.AcceptanceRate, &university.Type, &university.Category,
		&university.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return university, nil
}

func (r *PostgresUniversityRepository) SaveUniversity(ctx context.Context, savedUniversity *entity.SavedUniversity) error {
	query := `
		INSERT INTO saved_universities (id, user_id, university_id, saved_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query,
		savedUniversity.ID, savedUniversity.UserID,
		savedUniversity.UniversityID, savedUniversity.SavedAt,
	)
	return err
}

func (r *PostgresUniversityRepository) GetSavedUniversities(ctx context.Context, userID string) ([]*entity.University, error) {
	query := `
		SELECT u.id, u.name, u.country, u.acceptance_rate, u.type, u.category, u.created_at
		FROM universities u
		INNER JOIN saved_universities su ON su.university_id = u.id
		WHERE su.user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var universities []*entity.University
	for rows.Next() {
		university := &entity.University{}
		if err := rows.Scan(
			&university.ID, &university.Name, &university.Country,
			&university.AcceptanceRate, &university.Type, &university.Category,
			&university.CreatedAt,
		); err != nil {
			return nil, err
		}
		universities = append(universities, university)
	}
	return universities, nil
}

func (r *PostgresUniversityRepository) SearchGrants(ctx context.Context, country string, gpa float32) ([]*entity.Grant, error) {
	query := `
		SELECT id, name, description, amount, deadline, country, min_gpa, created_at
		FROM grants
		WHERE ($1 = '' OR country = $1)
		AND min_gpa <= $2
	`
	rows, err := r.db.QueryContext(ctx, query, country, gpa)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grants []*entity.Grant
	for rows.Next() {
		grant := &entity.Grant{}
		if err := rows.Scan(
			&grant.ID, &grant.Name, &grant.Description, &grant.Amount,
			&grant.Deadline, &grant.Country, &grant.MinGPA, &grant.CreatedAt,
		); err != nil {
			return nil, err
		}
		grants = append(grants, grant)
	}
	return grants, nil
}

func (r *PostgresUniversityRepository) GetGrantByID(ctx context.Context, grantID string) (*entity.Grant, error) {
	query := `SELECT id, name, description, amount, deadline, country, min_gpa, created_at FROM grants WHERE id = $1`
	grant := &entity.Grant{}
	err := r.db.QueryRowContext(ctx, query, grantID).Scan(
		&grant.ID, &grant.Name, &grant.Description, &grant.Amount,
		&grant.Deadline, &grant.Country, &grant.MinGPA, &grant.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return grant, nil
}

func (r *PostgresUniversityRepository) SaveGrant(ctx context.Context, savedGrant *entity.SavedGrant) error {
	query := `
		INSERT INTO saved_grants (id, user_id, grant_id, saved_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query,
		savedGrant.ID, savedGrant.UserID,
		savedGrant.GrantID, savedGrant.SavedAt,
	)
	return err
}

func (r *PostgresUniversityRepository) GetSavedGrants(ctx context.Context, userID string) ([]*entity.Grant, error) {
	query := `
		SELECT g.id, g.name, g.description, g.amount, g.deadline, g.country, g.min_gpa, g.created_at
		FROM grants g
		INNER JOIN saved_grants sg ON sg.grant_id = g.id
		WHERE sg.user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grants []*entity.Grant
	for rows.Next() {
		grant := &entity.Grant{}
		if err := rows.Scan(
			&grant.ID, &grant.Name, &grant.Description, &grant.Amount,
			&grant.Deadline, &grant.Country, &grant.MinGPA, &grant.CreatedAt,
		); err != nil {
			return nil, err
		}
		grants = append(grants, grant)
	}
	return grants, nil
}

func (r *PostgresUniversityRepository) saveEntity(ctx context.Context, id, userID, entityID, tableName string, savedAt time.Time) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, user_id, grant_id, saved_at) VALUES ($1, $2, $3, $4)`, tableName)
	_, err := r.db.ExecContext(ctx, query, id, userID, entityID, savedAt)
	return err
}

func (r *PostgresUniversityRepository) RemoveSavedUniversity(ctx context.Context, userID, universityID string) error {
	query := `DELETE FROM saved_universities WHERE user_id = $1 AND university_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, universityID)
	return err
}

func (r *PostgresUniversityRepository) RemoveSavedGrant(ctx context.Context, userID, grantID string) error {
	query := `DELETE FROM saved_grants WHERE user_id = $1 AND grant_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, grantID)
	return err
}

func (r *PostgresUniversityRepository) ListUniversitiesByCategory(ctx context.Context, category string) ([]*entity.University, error) {
	query := `
    SELECT id, name, country, acceptance_rate, type, category, created_at
    FROM universities WHERE category = $1
  `
	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var universities []*entity.University
	for rows.Next() {
		university := &entity.University{}
		if err := rows.Scan(
			&university.ID, &university.Name, &university.Country,
			&university.AcceptanceRate, &university.Type, &university.Category,
			&university.CreatedAt,
		); err != nil {
			return nil, err
		}
		universities = append(universities, university)
	}
	return universities, nil
}

func (r *PostgresUniversityRepository) GetUniversityStatistics(ctx context.Context) (int32, int32, int32, int32, error) {
	var total, reach, target, safety int32

	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM universities`).Scan(&total)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM universities WHERE category = 'reach'`).Scan(&reach)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM universities WHERE category = 'target'`).Scan(&target)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM universities WHERE category = 'safety'`).Scan(&safety)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return total, reach, target, safety, nil
}
