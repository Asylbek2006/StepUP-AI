package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/stepup-ai/university-service/internal/entity"
	"github.com/stepup-ai/university-service/internal/usecase"
)

type MockUniversityRepository struct {
	mock.Mock
}

func (m *MockUniversityRepository) SearchUniversities(ctx context.Context, country, universityType string, minAcceptanceRate, maxAcceptanceRate float32) ([]*entity.University, error) {
	args := m.Called(ctx, country, universityType, minAcceptanceRate, maxAcceptanceRate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.University), args.Error(1)
}

func (m *MockUniversityRepository) GetUniversityByID(ctx context.Context, universityID string) (*entity.University, error) {
	args := m.Called(ctx, universityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.University), args.Error(1)
}

func (m *MockUniversityRepository) SaveUniversity(ctx context.Context, savedUniversity *entity.SavedUniversity) error {
	args := m.Called(ctx, savedUniversity)
	return args.Error(0)
}

func (m *MockUniversityRepository) GetSavedUniversities(ctx context.Context, userID string) ([]*entity.University, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.University), args.Error(1)
}

func (m *MockUniversityRepository) SearchGrants(ctx context.Context, country string, gpa float32) ([]*entity.Grant, error) {
	args := m.Called(ctx, country, gpa)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Grant), args.Error(1)
}

func (m *MockUniversityRepository) GetGrantByID(ctx context.Context, grantID string) (*entity.Grant, error) {
	args := m.Called(ctx, grantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Grant), args.Error(1)
}

func (m *MockUniversityRepository) SaveGrant(ctx context.Context, savedGrant *entity.SavedGrant) error {
	args := m.Called(ctx, savedGrant)
	return args.Error(0)
}

func (m *MockUniversityRepository) GetSavedGrants(ctx context.Context, userID string) ([]*entity.Grant, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Grant), args.Error(1)
}

func TestSearchUniversities_Success(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	expectedUniversities := []*entity.University{
		{ID: "uni-1", Name: "MIT", Country: "USA", AcceptanceRate: 7.3},
		{ID: "uni-2", Name: "Harvard", Country: "USA", AcceptanceRate: 5.0},
	}

	mockRepository.On("SearchUniversities", mock.Anything, "USA", "", float32(0), float32(0)).Return(expectedUniversities, nil)

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	universities, err := universityUsecase.SearchUniversities(context.Background(), "USA", "", 0, 0)

	assert.NoError(t, err)
	assert.Len(t, universities, 2)
	assert.Equal(t, "MIT", universities[0].Name)
	mockRepository.AssertExpectations(t)
}

func TestSearchUniversities_EmptyResult(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	mockRepository.On("SearchUniversities", mock.Anything, "Kazakhstan", "", float32(0), float32(0)).Return([]*entity.University{}, nil)

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	universities, err := universityUsecase.SearchUniversities(context.Background(), "Kazakhstan", "", 0, 0)

	assert.NoError(t, err)
	assert.Empty(t, universities)
	mockRepository.AssertExpectations(t)
}

func TestGetUniversityDetails_Success(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	expectedUniversity := &entity.University{
		ID:             "uni-1",
		Name:           "MIT",
		Country:        "USA",
		AcceptanceRate: 7.3,
		Type:           "private",
		Category:       "reach",
	}

	mockRepository.On("GetUniversityByID", mock.Anything, "uni-1").Return(expectedUniversity, nil)

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	university, err := universityUsecase.GetUniversityDetails(context.Background(), "uni-1")

	assert.NoError(t, err)
	assert.Equal(t, "MIT", university.Name)
	assert.Equal(t, "USA", university.Country)
	mockRepository.AssertExpectations(t)
}

func TestSaveUniversity_Success(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	mockRepository.On("SaveUniversity", mock.Anything, mock.AnythingOfType("*entity.SavedUniversity")).Return(nil)

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	err := universityUsecase.SaveUniversity(context.Background(), "user-1", "uni-1")

	assert.NoError(t, err)
	mockRepository.AssertExpectations(t)
}

func TestGetSavedUniversities_Success(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	expectedUniversities := []*entity.University{
		{ID: "uni-1", Name: "MIT", Country: "USA"},
	}

	mockRepository.On("GetSavedUniversities", mock.Anything, "user-1").Return(expectedUniversities, nil)

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	universities, err := universityUsecase.GetSavedUniversities(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Len(t, universities, 1)
	mockRepository.AssertExpectations(t)
}

func TestSearchGrants_Success(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	expectedGrants := []*entity.Grant{
		{ID: "grant-1", Name: "Bolashak", Country: "Kazakhstan", MinGPA: 3.5},
		{ID: "grant-2", Name: "Chevening", Country: "UK", MinGPA: 3.0},
	}

	mockRepository.On("SearchGrants", mock.Anything, "Kazakhstan", float32(3.8)).Return(expectedGrants, nil)

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	grants, err := universityUsecase.SearchGrants(context.Background(), "Kazakhstan", 3.8)

	assert.NoError(t, err)
	assert.Len(t, grants, 2)
	assert.Equal(t, "Bolashak", grants[0].Name)
	mockRepository.AssertExpectations(t)
}

func TestSaveGrant_Success(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	mockRepository.On("SaveGrant", mock.Anything, mock.AnythingOfType("*entity.SavedGrant")).Return(nil)

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	err := universityUsecase.SaveGrant(context.Background(), "user-1", "grant-1")

	assert.NoError(t, err)
	mockRepository.AssertExpectations(t)
}

func TestGetSavedGrants_Success(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	expectedGrants := []*entity.Grant{
		{ID: "grant-1", Name: "Bolashak", Country: "Kazakhstan"},
	}

	mockRepository.On("GetSavedGrants", mock.Anything, "user-1").Return(expectedGrants, nil)

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	grants, err := universityUsecase.GetSavedGrants(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Len(t, grants, 1)
	assert.Equal(t, "Bolashak", grants[0].Name)
	mockRepository.AssertExpectations(t)
}
