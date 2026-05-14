package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/stepup-ai/api-gateway/proto/ai"
)

type AIHandler struct {
	aiServiceClient pb.AIServiceClient
}

func NewAIHandler(aiServiceConnection *grpc.ClientConn) *AIHandler {
	return &AIHandler{
		aiServiceClient: pb.NewAIServiceClient(aiServiceConnection),
	}
}

func (h *AIHandler) AnalyzeAdmissionChances(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var requestBody struct {
		UniversityID string  `json:"university_id"`
		GPA          float32 `json:"gpa"`
		SATScore     int32   `json:"sat_score"`
		IELTSScore   float32 `json:"ielts_score"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.aiServiceClient.AnalyzeAdmissionChances(ctx, &pb.AnalyzeAdmissionChancesRequest{
		UserId:       userID,
		UniversityId: requestBody.UniversityID,
		Gpa:          requestBody.GPA,
		SatScore:     requestBody.SATScore,
		IeltsScore:   requestBody.IELTSScore,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to analyze admission chances"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"admission_chance": response.AdmissionChancePercentage,
		"weak_areas":       response.WeakAreas,
		"suggestions":      response.ImprovementSuggestions,
		"gap_analysis":     response.GapAnalysis,
	})
}

func (h *AIHandler) GenerateRoadmap(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var requestBody struct {
		TargetUniversityID     string `json:"target_university_id"`
		MonthsUntilApplication int32  `json:"months_until_application"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.aiServiceClient.GenerateRoadmap(ctx, &pb.GenerateRoadmapRequest{
		UserId:                 userID,
		TargetUniversityId:     requestBody.TargetUniversityID,
		MonthsUntilApplication: requestBody.MonthsUntilApplication,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate roadmap"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"steps": response.Steps})
}

func (h *AIHandler) ReviewEssay(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var requestBody struct {
		EssayText      string `json:"essay_text"`
		UniversityName string `json:"university_name"`
		ProgramName    string `json:"program_name"`
		WordLimit      int32  `json:"word_limit"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.aiServiceClient.ReviewEssay(ctx, &pb.ReviewEssayRequest{
		UserId:         userID,
		EssayText:      requestBody.EssayText,
		UniversityName: requestBody.UniversityName,
		ProgramName:    requestBody.ProgramName,
		WordLimit:      requestBody.WordLimit,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to review essay"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"grammar_score":    response.GrammarScore,
		"coherence_score":  response.CoherenceScore,
		"uniqueness_score": response.UniquenessScore,
		"relevance_score":  response.RelevanceScore,
		"suggestions":      response.ImprovementSuggestions,
	})
}

func (h *AIHandler) MatchGrants(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var requestBody struct {
		GPA          float32  `json:"gpa"`
		Country      string   `json:"country"`
		Achievements []string `json:"achievements"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.aiServiceClient.MatchGrants(ctx, &pb.MatchGrantsRequest{
		UserId:       userID,
		Gpa:          requestBody.GPA,
		Country:      requestBody.Country,
		Achievements: requestBody.Achievements,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to match grants"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"matched_grant_ids": response.MatchedGrantIds})
}

func (h *AIHandler) GetAnalysisHistory(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	response, err := h.aiServiceClient.GetAnalysisHistory(ctx, &pb.GetAnalysisHistoryRequest{
		UserId: userID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get analysis history"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"analysis_ids": response.AnalysisIds})
}
