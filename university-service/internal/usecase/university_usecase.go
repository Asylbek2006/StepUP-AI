package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/stepup-ai/university-service/internal/entity"
	"github.com/stepup-ai/university-service/internal/repository"
)

type UniversityUsecase interface {
	SearchUniversities(ctx context.Context, country, universityType string, minAcceptanceRate, maxAcceptanceRate float32) ([]*entity.University, error)
	GetUniversityDetails(ctx context.Context, universityID string) (*entity.University, error)
	SaveUniversity(ctx context.Context, userID, universityID string) error
	GetSavedUniversities(ctx context.Context, userID string) ([]*entity.University, error)
	SearchGrants(ctx context.Context, country string, gpa float32) ([]*entity.Grant, error)
	SaveGrant(ctx context.Context, userID, grantID string) error
	GetSavedGrants(ctx context.Context, userID string) ([]*entity.Grant, error)
	RemoveSavedUniversity(ctx context.Context, userID, universityID string) error
	RemoveSavedGrant(ctx context.Context, userID, grantID string) error
	GetGrantDetails(ctx context.Context, grantID string) (*entity.Grant, error)
	ListUniversitiesByCategory(ctx context.Context, category string) ([]*entity.University, error)
	GetUniversityStatistics(ctx context.Context) (int32, int32, int32, int32, error)
}

type universityUsecase struct {
	universityRepository repository.UniversityRepository
}

func NewUniversityUsecase(universityRepository repository.UniversityRepository) UniversityUsecase {
	return &universityUsecase{universityRepository: universityRepository}
}

func (u *universityUsecase) SearchUniversities(ctx context.Context, country, universityType string, minAcceptanceRate, maxAcceptanceRate float32) ([]*entity.University, error) {
	return u.universityRepository.SearchUniversities(ctx, country, universityType, minAcceptanceRate, maxAcceptanceRate)
}

func (u *universityUsecase) GetUniversityDetails(ctx context.Context, universityID string) (*entity.University, error) {
	return u.universityRepository.GetUniversityByID(ctx, universityID)
}

func (u *universityUsecase) SaveUniversity(ctx context.Context, userID, universityID string) error {
	savedUniversity := &entity.SavedUniversity{
		ID:           uuid.New().String(),
		UserID:       userID,
		UniversityID: universityID,
		SavedAt:      time.Now(),
	}
	return u.universityRepository.SaveUniversity(ctx, savedUniversity)
}

func (u *universityUsecase) GetSavedUniversities(ctx context.Context, userID string) ([]*entity.University, error) {
	return u.universityRepository.GetSavedUniversities(ctx, userID)
}

func (u *universityUsecase) SearchGrants(ctx context.Context, country string, gpa float32) ([]*entity.Grant, error) {
	return u.universityRepository.SearchGrants(ctx, country, gpa)
}

func (u *universityUsecase) SaveGrant(ctx context.Context, userID, grantID string) error {
	savedGrant := &entity.SavedGrant{
		ID:      uuid.New().String(),
		UserID:  userID,
		GrantID: grantID,
		SavedAt: time.Now(),
	}
	return u.universityRepository.SaveGrant(ctx, savedGrant)
}

func (u *universityUsecase) GetSavedGrants(ctx context.Context, userID string) ([]*entity.Grant, error) {
	return u.universityRepository.GetSavedGrants(ctx, userID)
}

func (u *universityUsecase) RemoveSavedUniversity(ctx context.Context, userID, universityID string) error {
	return u.universityRepository.RemoveSavedUniversity(ctx, userID, universityID)
}

func (u *universityUsecase) RemoveSavedGrant(ctx context.Context, userID, grantID string) error {
	return u.universityRepository.RemoveSavedGrant(ctx, userID, grantID)
}

func (u *universityUsecase) GetGrantDetails(ctx context.Context, grantID string) (*entity.Grant, error) {
	return u.universityRepository.GetGrantByID(ctx, grantID)
}

func (u *universityUsecase) ListUniversitiesByCategory(ctx context.Context, category string) ([]*entity.University, error) {
	return u.universityRepository.ListUniversitiesByCategory(ctx, category)
}

func (u *universityUsecase) GetUniversityStatistics(ctx context.Context) (int32, int32, int32, int32, error) {
	return u.universityRepository.GetUniversityStatistics(ctx)
}
