package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/stepup-ai/user-service/internal/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, userID string) (*entity.User, error)
	CreateUserProfile(ctx context.Context, profile *entity.UserProfile) error
	GetUserProfile(ctx context.Context, userID string) (*entity.UserProfile, error)
	UpdateUserProfile(ctx context.Context, profile *entity.UserProfile) error
	SaveRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	SavePasswordResetToken(ctx context.Context, token *entity.PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error)
	MarkPasswordResetTokenAsUsed(ctx context.Context, tokenHash string) error
	UpdateUserPassword(ctx context.Context, userID string, passwordHash string) error
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash,
		user.FullName, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, email, password_hash, full_name, created_at, updated_at FROM users WHERE email = $1`
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.FullName, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	query := `SELECT id, email, password_hash, full_name, created_at, updated_at FROM users WHERE id = $1`
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.FullName, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) CreateUserProfile(ctx context.Context, profile *entity.UserProfile) error {
	query := `
		INSERT INTO user_profiles (user_id, gpa, sat_score, ielts_score, toefl_score, country, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		profile.UserID, profile.GPA, profile.SATScore,
		profile.IELTSScore, profile.TOEFLScore, profile.Country, profile.UpdatedAt,
	)
	return err
}

func (r *PostgresUserRepository) GetUserProfile(ctx context.Context, userID string) (*entity.UserProfile, error) {
	query := `SELECT user_id, gpa, sat_score, ielts_score, toefl_score, country, updated_at FROM user_profiles WHERE user_id = $1`
	profile := &entity.UserProfile{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.UserID, &profile.GPA, &profile.SATScore,
		&profile.IELTSScore, &profile.TOEFLScore, &profile.Country, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *PostgresUserRepository) UpdateUserProfile(ctx context.Context, profile *entity.UserProfile) error {
	query := `
		UPDATE user_profiles
		SET gpa = $1, sat_score = $2, ielts_score = $3, toefl_score = $4, country = $5, updated_at = $6
		WHERE user_id = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		profile.GPA, profile.SATScore, profile.IELTSScore,
		profile.TOEFLScore, profile.Country, time.Now(), profile.UserID,
	)
	return err
}

func (r *PostgresUserRepository) SaveRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt,
	)
	return err
}

func (r *PostgresUserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (r *PostgresUserRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, created_at FROM refresh_tokens WHERE token_hash = $1`
	token := &entity.RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *PostgresUserRepository) SavePasswordResetToken(ctx context.Context, token *entity.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at, used)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.Used,
	)
	return err
}

func (r *PostgresUserRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, used FROM password_reset_tokens WHERE token_hash = $1`
	token := &entity.PasswordResetToken{}
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.Used,
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *PostgresUserRepository) MarkPasswordResetTokenAsUsed(ctx context.Context, tokenHash string) error {
	query := `UPDATE password_reset_tokens SET used = TRUE WHERE token_hash = $1`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (r *PostgresUserRepository) UpdateUserPassword(ctx context.Context, userID string, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, passwordHash, time.Now(), userID)
	return err
}
