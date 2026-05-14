package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/stepup-ai/user-service/internal/entity"
	"github.com/stepup-ai/user-service/internal/messaging"
	"github.com/stepup-ai/user-service/internal/repository"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrResetTokenInvalid  = errors.New("reset token is invalid or expired")
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, email, password, fullName string) (accessToken, refreshToken string, err error)
	LoginUser(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	LogoutUser(ctx context.Context, refreshToken string) error
	RefreshAccessToken(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	GetUserProfile(ctx context.Context, userID string) (*entity.UserProfile, error)
	UpdateUserProfile(ctx context.Context, profile *entity.UserProfile) error
	SendPasswordResetEmail(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, resetToken, newPassword string) error
}

type userUsecase struct {
	userRepository repository.UserRepository
	jwtSecretKey   string
	emailSender    EmailSender
	natsPublisher  *messaging.NatsPublisher
}

type EmailSender interface {
	SendPasswordResetEmail(toEmail, resetLink string) error
}

func NewUserUsecase(
	userRepository repository.UserRepository,
	jwtSecretKey string,
	emailSender EmailSender,
	natsPublisher *messaging.NatsPublisher,
) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
		jwtSecretKey:   jwtSecretKey,
		emailSender:    emailSender,
		natsPublisher:  natsPublisher,
	}
}

func (u *userUsecase) RegisterUser(ctx context.Context, email, password, fullName string) (string, string, error) {
	existingUser, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return "", "", err
	}
	if existingUser != nil {
		return "", "", ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	newUser := &entity.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(passwordHash),
		FullName:     fullName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.userRepository.CreateUser(ctx, newUser); err != nil {
		return "", "", err
	}

	u.natsPublisher.PublishUserRegisteredEvent(messaging.UserRegisteredEvent{
		UserID:   newUser.ID,
		Email:    newUser.Email,
		FullName: newUser.FullName,
	})

	emptyProfile := &entity.UserProfile{
		UserID:    newUser.ID,
		UpdatedAt: time.Now(),
	}
	if err := u.userRepository.CreateUserProfile(ctx, emptyProfile); err != nil {
		return "", "", err
	}

	accessToken, err := u.generateAccessToken(newUser.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := u.createAndSaveRefreshToken(ctx, newUser.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *userUsecase) LoginUser(ctx context.Context, email, password string) (string, string, error) {
	user, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := u.generateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := u.createAndSaveRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *userUsecase) LogoutUser(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	return u.userRepository.DeleteRefreshToken(ctx, tokenHash)
}

func (u *userUsecase) RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
	tokenHash := hashToken(refreshToken)

	savedToken, err := u.userRepository.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	if time.Now().After(savedToken.ExpiresAt) {
		return "", "", ErrTokenExpired
	}

	if err := u.userRepository.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return "", "", err
	}

	newAccessToken, err := u.generateAccessToken(savedToken.UserID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := u.createAndSaveRefreshToken(ctx, savedToken.UserID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (u *userUsecase) GetUserProfile(ctx context.Context, userID string) (*entity.UserProfile, error) {
	return u.userRepository.GetUserProfile(ctx, userID)
}

func (u *userUsecase) UpdateUserProfile(ctx context.Context, profile *entity.UserProfile) error {
	return u.userRepository.UpdateUserProfile(ctx, profile)
}

func (u *userUsecase) SendPasswordResetEmail(ctx context.Context, email string) error {
	user, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	resetToken := generateSecureToken()
	tokenHash := hashToken(resetToken)

	passwordResetToken := &entity.PasswordResetToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	if err := u.userRepository.SavePasswordResetToken(ctx, passwordResetToken); err != nil {
		return err
	}

	resetLink := "https://stepupai.com/reset-password?token=" + resetToken
	return u.emailSender.SendPasswordResetEmail(email, resetLink)
}

func (u *userUsecase) ResetPassword(ctx context.Context, resetToken, newPassword string) error {
	tokenHash := hashToken(resetToken)

	savedToken, err := u.userRepository.GetPasswordResetToken(ctx, tokenHash)
	if err != nil {
		return ErrResetTokenInvalid
	}

	if savedToken.Used || time.Now().After(savedToken.ExpiresAt) {
		return ErrResetTokenInvalid
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := u.userRepository.UpdateUserPassword(ctx, savedToken.UserID, string(newPasswordHash)); err != nil {
		return err
	}

	return u.userRepository.MarkPasswordResetTokenAsUsed(ctx, tokenHash)
}

func (u *userUsecase) generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtSecretKey))
}

func (u *userUsecase) createAndSaveRefreshToken(ctx context.Context, userID string) (string, error) {
	refreshToken := generateSecureToken()
	tokenHash := hashToken(refreshToken)

	savedRefreshToken := &entity.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := u.userRepository.SaveRefreshToken(ctx, savedRefreshToken); err != nil {
		return "", err
	}

	return refreshToken, nil
}

func generateSecureToken() string {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	return hex.EncodeToString(randomBytes)
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
