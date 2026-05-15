package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/stepup-ai/ai-service/internal/entity"
)

type MockAIRepository struct {
	mock.Mock
}

func (m *MockAIRepository) SaveAdmissionAnalysis(ctx context.Context, analysis *entity.AdmissionAnalysis) error {
	args := m.Called(ctx, analysis)
	return args.Error(0)
}

func (m *MockAIRepository) GetAdmissionAnalysisByID(ctx context.Context, analysisID string) (*entity.AdmissionAnalysis, error) {
	args := m.Called(ctx, analysisID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AdmissionAnalysis), args.Error(1)
}

func (m *MockAIRepository) GetAdmissionAnalysisHistoryByUserID(ctx context.Context, userID string) ([]*entity.AdmissionAnalysis, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.AdmissionAnalysis), args.Error(1)
}

func (m *MockAIRepository) SaveRoadmap(ctx context.Context, roadmap *entity.Roadmap, steps []*entity.RoadmapStep) error {
	args := m.Called(ctx, roadmap, steps)
	return args.Error(0)
}

func (m *MockAIRepository) GetRoadmapByUserID(ctx context.Context, userID string) (*entity.Roadmap, []*entity.RoadmapStep, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*entity.Roadmap), args.Get(1).([]*entity.RoadmapStep), args.Error(2)
}

func (m *MockAIRepository) SaveEssayReview(ctx context.Context, review *entity.EssayReview) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func (m *MockAIRepository) GetEssayReviewsByUserID(ctx context.Context, userID string) ([]*entity.EssayReview, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EssayReview), args.Error(1)
}

func TestMockRepository_SaveAdmissionAnalysis(t *testing.T) {
	mockRepository := new(MockAIRepository)

	analysis := &entity.AdmissionAnalysis{
		ID:                        "analysis-1",
		UserID:                    "user-1",
		UniversityID:              "uni-1",
		AdmissionChancePercentage: 75.5,
		WeakAreas:                 []string{"SAT", "Essays"},
		ImprovementSuggestions:    []string{"Improve SAT score", "Write better essays"},
		GapAnalysis:               "Need to improve test scores",
		CreatedAt:                 time.Now(),
	}

	mockRepository.On("SaveAdmissionAnalysis", mock.Anything, analysis).Return(nil)

	err := mockRepository.SaveAdmissionAnalysis(context.Background(), analysis)

	assert.NoError(t, err)
	mockRepository.AssertExpectations(t)
}

func TestMockRepository_GetAdmissionAnalysisByID(t *testing.T) {
	mockRepository := new(MockAIRepository)

	expectedAnalysis := &entity.AdmissionAnalysis{
		ID:                        "analysis-1",
		UserID:                    "user-1",
		UniversityID:              "uni-1",
		AdmissionChancePercentage: 80.0,
	}

	mockRepository.On("GetAdmissionAnalysisByID", mock.Anything, "analysis-1").Return(expectedAnalysis, nil)

	analysis, err := mockRepository.GetAdmissionAnalysisByID(context.Background(), "analysis-1")

	assert.NoError(t, err)
	assert.Equal(t, "analysis-1", analysis.ID)
	assert.Equal(t, float32(80.0), analysis.AdmissionChancePercentage)
	mockRepository.AssertExpectations(t)
}

func TestMockRepository_GetAdmissionAnalysisHistoryByUserID(t *testing.T) {
	mockRepository := new(MockAIRepository)

	expectedAnalyses := []*entity.AdmissionAnalysis{
		{ID: "analysis-1", UserID: "user-1"},
		{ID: "analysis-2", UserID: "user-1"},
	}

	mockRepository.On("GetAdmissionAnalysisHistoryByUserID", mock.Anything, "user-1").Return(expectedAnalyses, nil)

	analyses, err := mockRepository.GetAdmissionAnalysisHistoryByUserID(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Len(t, analyses, 2)
	mockRepository.AssertExpectations(t)
}

func TestMockRepository_SaveRoadmap(t *testing.T) {
	mockRepository := new(MockAIRepository)

	roadmap := &entity.Roadmap{
		ID:        "roadmap-1",
		UserID:    "user-1",
		CreatedAt: time.Now(),
	}

	steps := []*entity.RoadmapStep{
		{ID: "step-1", RoadmapID: "roadmap-1", Title: "IELTS Prep", MonthNumber: 1},
		{ID: "step-2", RoadmapID: "roadmap-1", Title: "SAT Prep", MonthNumber: 2},
	}

	mockRepository.On("SaveRoadmap", mock.Anything, roadmap, steps).Return(nil)

	err := mockRepository.SaveRoadmap(context.Background(), roadmap, steps)

	assert.NoError(t, err)
	mockRepository.AssertExpectations(t)
}

func TestMockRepository_SaveEssayReview(t *testing.T) {
	mockRepository := new(MockAIRepository)

	review := &entity.EssayReview{
		ID:              "review-1",
		UserID:          "user-1",
		EssayText:       "My essay text",
		UniversityName:  "MIT",
		ProgramName:     "Computer Science",
		GrammarScore:    8.5,
		CoherenceScore:  7.0,
		UniquenessScore: 9.0,
		RelevanceScore:  8.0,
		CreatedAt:       time.Now(),
	}

	mockRepository.On("SaveEssayReview", mock.Anything, review).Return(nil)

	err := mockRepository.SaveEssayReview(context.Background(), review)

	assert.NoError(t, err)
	mockRepository.AssertExpectations(t)
}

func TestMockRepository_GetEssayReviewsByUserID(t *testing.T) {
	mockRepository := new(MockAIRepository)

	expectedReviews := []*entity.EssayReview{
		{ID: "review-1", UserID: "user-1", UniversityName: "MIT"},
		{ID: "review-2", UserID: "user-1", UniversityName: "Harvard"},
	}

	mockRepository.On("GetEssayReviewsByUserID", mock.Anything, "user-1").Return(expectedReviews, nil)

	reviews, err := mockRepository.GetEssayReviewsByUserID(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Len(t, reviews, 2)
	mockRepository.AssertExpectations(t)
}
