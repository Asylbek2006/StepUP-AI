package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/stepup-ai/university-service/internal/entity"
	"github.com/stepup-ai/university-service/internal/usecase"
)

func TestGetUniversityDetails_NotFound(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	mockRepository.On("GetUniversityByID", mock.Anything, "non-existent").Return(nil, errors.New("not found"))

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	university, err := universityUsecase.GetUniversityDetails(context.Background(), "non-existent")

	assert.Error(t, err)
	assert.Nil(t, university)
	mockRepository.AssertExpectations(t)
}

func TestSearchUniversities_WithFilters(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	expectedUniversities := []*entity.University{
		{ID: "uni-1", Name: "MIT", Country: "USA", AcceptanceRate: 6.7, Type: "private"},
	}

	mockRepository.On("SearchUniversities", mock.Anything, "USA", "private", float32(0), float32(10)).Return(expectedUniversities, nil)

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	universities, err := universityUsecase.SearchUniversities(context.Background(), "USA", "private", 0, 10)

	assert.NoError(t, err)
	assert.Len(t, universities, 1)
	assert.Equal(t, "MIT", universities[0].Name)
	mockRepository.AssertExpectations(t)
}

func TestSaveGrant_RepositoryError(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	mockRepository.On("SaveGrant", mock.Anything, mock.AnythingOfType("*entity.SavedGrant")).Return(errors.New("database error"))

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	err := universityUsecase.SaveGrant(context.Background(), "user-1", "grant-1")

	assert.Error(t, err)
	mockRepository.AssertExpectations(t)
}

func TestGetSavedUniversities_EmptyResult(t *testing.T) {
	mockRepository := new(MockUniversityRepository)

	mockRepository.On("GetSavedUniversities", mock.Anything, "user-without-saved").Return([]*entity.University{}, nil)

	universityUsecase := usecase.NewUniversityUsecase(mockRepository)

	universities, err := universityUsecase.GetSavedUniversities(context.Background(), "user-without-saved")

	assert.NoError(t, err)
	assert.Empty(t, universities)
	mockRepository.AssertExpectations(t)
}
