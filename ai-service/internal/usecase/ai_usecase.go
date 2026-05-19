package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"

	"github.com/stepup-ai/ai-service/internal/entity"
	"github.com/stepup-ai/ai-service/internal/repository"
)

type AIUsecase interface {
	AnalyzeAdmissionChances(ctx context.Context, userID, universityID string, gpa float32, satScore int32, ieltsScore float32) (*entity.AdmissionAnalysis, error)
	GenerateRoadmap(ctx context.Context, userID, targetUniversityID string, monthsUntilApplication int32) (*entity.Roadmap, []*entity.RoadmapStep, error)
	ReviewEssay(ctx context.Context, userID, essayText, universityName, programName string, wordLimit int32) (*entity.EssayReview, error)
	MatchGrants(ctx context.Context, userID, country string, gpa float32, achievements []string) ([]string, error)
	GetAnalysisHistory(ctx context.Context, userID string) ([]*entity.AdmissionAnalysis, error)
}

type aiUsecase struct {
	aiRepository repository.AIRepository
	openaiClient *openai.Client
}

func NewAIUsecase(aiRepository repository.AIRepository, openaiAPIKey string) AIUsecase {
	return &aiUsecase{
		aiRepository: aiRepository,
		openaiClient: openai.NewClient(openaiAPIKey),
	}
}

func (u *aiUsecase) AnalyzeAdmissionChances(ctx context.Context, userID, universityID string, gpa float32, satScore int32, ieltsScore float32) (*entity.AdmissionAnalysis, error) {
	prompt := fmt.Sprintf(`
    Analyze admission chances for a student with the following profile:
    GPA: %.2f
    SAT Score: %d
    IELTS Score: %.1f
    University ID: %s

    Respond ONLY in this JSON format:
    {
      "admission_chance_percentage": 75.5,
      "weak_areas": ["area1", "area2"],
      "improvement_suggestions": ["suggestion1", "suggestion2"],
      "gap_analysis": "detailed gap analysis text"
    }
  `, gpa, satScore, ieltsScore, universityID)

	response, err := u.openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		AdmissionChancePercentage float32  `json:"admission_chance_percentage"`
		WeakAreas                 []string `json:"weak_areas"`
		ImprovementSuggestions    []string `json:"improvement_suggestions"`
		GapAnalysis               string   `json:"gap_analysis"`
	}

	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &result); err != nil {
		return nil, err
	}

	analysis := &entity.AdmissionAnalysis{
		ID:                        uuid.New().String(),
		UserID:                    userID,
		UniversityID:              universityID,
		AdmissionChancePercentage: result.AdmissionChancePercentage,
		WeakAreas:                 result.WeakAreas,
		ImprovementSuggestions:    result.ImprovementSuggestions,
		GapAnalysis:               result.GapAnalysis,
		CreatedAt:                 time.Now(),
	}

	if err := u.aiRepository.SaveAdmissionAnalysis(ctx, analysis); err != nil {
		return nil, err
	}

	return analysis, nil
}

func (u *aiUsecase) GenerateRoadmap(ctx context.Context, userID, targetUniversityID string, monthsUntilApplication int32) (*entity.Roadmap, []*entity.RoadmapStep, error) {
	prompt := fmt.Sprintf(`
    Generate a detailed preparation roadmap for a student applying to university.
    Months until application: %d
    Target University ID: %s

    Respond ONLY in this JSON format:
    {
      "steps": [
        {
          "title": "step title",
          "description": "detailed description",
          "month_number": 1,
          "category": "IELTS/SAT/Essay/Research/Application"
        }
      ]
    }
  `, monthsUntilApplication, targetUniversityID)
	response, err := u.openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Steps []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			MonthNumber int32  `json:"month_number"`
			Category    string `json:"category"`
		} `json:"steps"`
	}

	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &result); err != nil {
		return nil, nil, err
	}

	roadmap := &entity.Roadmap{
		ID:        uuid.New().String(),
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	var roadmapSteps []*entity.RoadmapStep
	for _, step := range result.Steps {
		roadmapSteps = append(roadmapSteps, &entity.RoadmapStep{
			ID:          uuid.New().String(),
			RoadmapID:   roadmap.ID,
			Title:       step.Title,
			Description: step.Description,
			MonthNumber: step.MonthNumber,
			Category:    step.Category,
		})
	}

	if err := u.aiRepository.SaveRoadmap(ctx, roadmap, roadmapSteps); err != nil {
		return nil, nil, err
	}

	return roadmap, roadmapSteps, nil
}

func (u *aiUsecase) ReviewEssay(ctx context.Context, userID, essayText, universityName, programName string, wordLimit int32) (*entity.EssayReview, error) {
	prompt := fmt.Sprintf(`
    Review the following university application essay:
    University: %s
    Program: %s
    Word Limit: %d
    Essay: %s

    Respond ONLY in this JSON format:
    {
      "grammar_score": 8.5,
      "coherence_score": 7.0,
      "uniqueness_score": 9.0,
      "relevance_score": 8.0,
      "improvement_suggestions": ["suggestion1", "suggestion2"]
    }
  `, universityName, programName, wordLimit, essayText)

	response, err := u.openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		GrammarScore           float32  `json:"grammar_score"`
		CoherenceScore         float32  `json:"coherence_score"`
		UniquenessScore        float32  `json:"uniqueness_score"`
		RelevanceScore         float32  `json:"relevance_score"`
		ImprovementSuggestions []string `json:"improvement_suggestions"`
	}

	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &result); err != nil {
		return nil, err
	}

	review := &entity.EssayReview{
		ID:                     uuid.New().String(),
		UserID:                 userID,
		EssayText:              essayText,
		UniversityName:         universityName,
		ProgramName:            programName,
		GrammarScore:           result.GrammarScore,
		CoherenceScore:         result.CoherenceScore,
		UniquenessScore:        result.UniquenessScore,
		RelevanceScore:         result.RelevanceScore,
		ImprovementSuggestions: result.ImprovementSuggestions,
		CreatedAt:              time.Now(),
	}

	if err := u.aiRepository.SaveEssayReview(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

func (u *aiUsecase) MatchGrants(ctx context.Context, userID, country string, gpa float32, achievements []string) ([]string, error) {
	prompt := fmt.Sprintf(`
    Match grants for a student with the following profile:
    Country: %s
    GPA: %.2f
    Achievements: %v

    Respond ONLY in this JSON format:
    {
      "matched_grant_ids": ["grant_id_1", "grant_id_2"]
    }
  `, country, gpa, achievements)

	response, err := u.openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		MatchedGrantIDs []string `json:"matched_grant_ids"`
	}
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &result); err != nil {
		return nil, err
	}

	return result.MatchedGrantIDs, nil
}

func (u *aiUsecase) GetAnalysisHistory(ctx context.Context, userID string) ([]*entity.AdmissionAnalysis, error) {
	return u.aiRepository.GetAdmissionAnalysisHistoryByUserID(ctx, userID)
}
