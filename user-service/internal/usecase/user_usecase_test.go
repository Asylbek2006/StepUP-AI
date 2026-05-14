package usecase_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/stepup-ai/user-service/internal/entity"
	"github.com/stepup-ai/user-service/internal/usecase"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) CreateUserProfile(ctx context.Context, profile *entity.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserProfile(ctx context.Context, userID string) (*entity.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserProfile), args.Error(1)
}

func (m *MockUserRepository) UpdateUserProfile(ctx context.Context, profile *entity.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockUserRepository) SaveRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *MockUserRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.RefreshToken), args.Error(1)
}

func (m *MockUserRepository) SavePasswordResetToken(ctx context.Context, token *entity.PasswordResetToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockUserRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.PasswordResetToken), args.Error(1)
}

func (m *MockUserRepository) MarkPasswordResetTokenAsUsed(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserPassword(ctx context.Context, userID string, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
	return args.Error(0)
}

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendPasswordResetEmail(toEmail, resetLink string) error {
	args := m.Called(toEmail, resetLink)
	return args.Error(0)
}

type MockNatsPublisher struct {
	mock.Mock
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepository := new(MockUserRepository)
	mockEmailSender := new(MockEmailSender)

	mockRepository.On("GetUserByEmail", mock.Anything, "test@gmail.com").Return(nil, sql.ErrNoRows)
	mockRepository.On("CreateUser", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
	mockRepository.On("CreateUserProfile", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
	mockRepository.On("SaveRefreshToken", mock.Anything, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)

	userUsecase := usecase.NewUserUsecase(mockRepository, "test-secret-key", mockEmailSender, nil)

	accessToken, refreshToken, err := userUsecase.RegisterUser(context.Background(), "test@gmail.com", "password123", "Test User")

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	mockRepository.AssertExpectations(t)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	mockRepository := new(MockUserRepository)
	mockEmailSender := new(MockEmailSender)

	existingUser := &entity.User{
		ID:    "existing-user-id",
		Email: "test@gmail.com",
	}

	mockRepository.On("GetUserByEmail", mock.Anything, "test@gmail.com").Return(existingUser, nil)

	userUsecase := usecase.NewUserUsecase(mockRepository, "test-secret-key", mockEmailSender, nil)

	_, _, err := userUsecase.RegisterUser(context.Background(), "test@gmail.com", "password123", "Test User")

	assert.ErrorIs(t, err, usecase.ErrUserAlreadyExists)
	mockRepository.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	mockRepository := new(MockUserRepository)
	mockEmailSender := new(MockEmailSender)

	userUsecase := usecase.NewUserUsecase(mockRepository, "test-secret-key", mockEmailSender, nil)

	mockRepository.On("GetUserByEmail", mock.Anything, "test@gmail.com").Return(nil, sql.ErrNoRows)
	mockRepository.On("CreateUser", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
	mockRepository.On("CreateUserProfile", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
	mockRepository.On("SaveRefreshToken", mock.Anything, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)

	_, _, _ = userUsecase.RegisterUser(context.Background(), "test@gmail.com", "password123", "Test User")

	mockRepository.On("GetUserByEmail", mock.Anything, "test@gmail.com").Return(&entity.User{
		ID:           "user-id",
		Email:        "test@gmail.com",
		PasswordHash: "$2a$10$test",
	}, nil)

	_, _, err := userUsecase.LoginUser(context.Background(), "test@gmail.com", "password123")
	assert.Error(t, err)
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	mockRepository := new(MockUserRepository)
	mockEmailSender := new(MockEmailSender)

	mockRepository.On("GetUserByEmail", mock.Anything, "test@gmail.com").Return(&entity.User{
		ID:           "user-id",
		Email:        "test@gmail.com",
		PasswordHash: "$2a$10$invalidhash",
	}, nil)

	userUsecase := usecase.NewUserUsecase(mockRepository, "test-secret-key", mockEmailSender, nil)

	_, _, err := userUsecase.LoginUser(context.Background(), "test@gmail.com", "wrongpassword")

	assert.ErrorIs(t, err, usecase.ErrInvalidCredentials)
	mockRepository.AssertExpectations(t)
}

func TestLogoutUser_Success(t *testing.T) {
	mockRepository := new(MockUserRepository)
	mockEmailSender := new(MockEmailSender)

	mockRepository.On("DeleteRefreshToken", mock.Anything, mock.AnythingOfType("string")).Return(nil)

	userUsecase := usecase.NewUserUsecase(mockRepository, "test-secret-key", mockEmailSender, nil)

	err := userUsecase.LogoutUser(context.Background(), "some-refresh-token")

	assert.NoError(t, err)
	mockRepository.AssertExpectations(t)
}

func TestRefreshAccessToken_InvalidToken(t *testing.T) {
	mockRepository := new(MockUserRepository)
	mockEmailSender := new(MockEmailSender)

	mockRepository.On("GetRefreshToken", mock.Anything, mock.AnythingOfType("string")).Return(nil, sql.ErrNoRows)

	userUsecase := usecase.NewUserUsecase(mockRepository, "test-secret-key", mockEmailSender, nil)

	_, _, err := userUsecase.RefreshAccessToken(context.Background(), "invalid-token")

	assert.ErrorIs(t, err, usecase.ErrInvalidToken)
	mockRepository.AssertExpectations(t)
}

func TestRefreshAccessToken_ExpiredToken(t *testing.T) {
	mockRepository := new(MockUserRepository)
	mockEmailSender := new(MockEmailSender)

	expiredToken := &entity.RefreshToken{
		ID:        "token-id",
		UserID:    "user-id",
		TokenHash: "hashed-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	mockRepository.On("GetRefreshToken", mock.Anything, mock.AnythingOfType("string")).Return(expiredToken, nil)

	userUsecase := usecase.NewUserUsecase(mockRepository, "test-secret-key", mockEmailSender, nil)

	_, _, err := userUsecase.RefreshAccessToken(context.Background(), "expired-token")

	assert.ErrorIs(t, err, usecase.ErrTokenExpired)
	mockRepository.AssertExpectations(t)
}
