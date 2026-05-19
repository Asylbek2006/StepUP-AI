package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stepup-ai/ai-service/internal/usecase"
	pb "github.com/stepup-ai/ai-service/proto/ai"
)

type AIGRPCHandler struct {
	pb.UnimplementedAIServiceServer
	aiUsecase usecase.AIUsecase
}

func NewAIGRPCHandler(u usecase.AIUsecase) *AIGRPCHandler {
	return &AIGRPCHandler{aiUsecase: u}
}

func (h *AIGRPCHandler) AnalyzeAdmissionChances(ctx context.Context, req *pb.AnalyzeAdmissionChancesRequest) (*pb.AnalyzeAdmissionChancesResponse, error) {
	analysis, err := h.aiUsecase.AnalyzeAdmissionChances(ctx, req.UserId, req.UniversityId, req.Gpa, req.SatScore, req.IeltsScore)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to analyze admission chances")
	}

	return &pb.AnalyzeAdmissionChancesResponse{
		AdmissionChancePercentage: analysis.AdmissionChancePercentage,
		WeakAreas:                 analysis.WeakAreas,
		ImprovementSuggestions:    analysis.ImprovementSuggestions,
		GapAnalysis:               analysis.GapAnalysis,
	}, nil
}

func (h *AIGRPCHandler) GenerateRoadmap(ctx context.Context, req *pb.GenerateRoadmapRequest) (*pb.GenerateRoadmapResponse, error) {
	_, steps, err := h.aiUsecase.GenerateRoadmap(ctx, req.UserId, req.TargetUniversityId, req.MonthsUntilApplication)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate roadmap")
	}

	var pbSteps []*pb.RoadmapStep
	for _, step := range steps {
		pbSteps = append(pbSteps, &pb.RoadmapStep{
			Title:       step.Title,
			Description: step.Description,
			MonthNumber: step.MonthNumber,
			Category:    step.Category,
		})
	}

	return &pb.GenerateRoadmapResponse{Steps: pbSteps}, nil
}

func (h *AIGRPCHandler) ReviewEssay(ctx context.Context, req *pb.ReviewEssayRequest) (*pb.ReviewEssayResponse, error) {
	review, err := h.aiUsecase.ReviewEssay(ctx, req.UserId, req.EssayText, req.UniversityName, req.ProgramName, req.WordLimit)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to review essay")
	}

	return &pb.ReviewEssayResponse{
		GrammarScore:           review.GrammarScore,
		CoherenceScore:         review.CoherenceScore,
		UniquenessScore:        review.UniquenessScore,
		RelevanceScore:         review.RelevanceScore,
		ImprovementSuggestions: review.ImprovementSuggestions,
	}, nil
}

func (h *AIGRPCHandler) MatchGrants(ctx context.Context, req *pb.MatchGrantsRequest) (*pb.MatchGrantsResponse, error) {
	matchedGrantIDs, err := h.aiUsecase.MatchGrants(ctx, req.UserId, req.Country, req.Gpa, req.Achievements)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to match grants")
	}

	return &pb.MatchGrantsResponse{MatchedGrantIds: matchedGrantIDs}, nil
}

func (h *AIGRPCHandler) GetAnalysisHistory(ctx context.Context, req *pb.GetAnalysisHistoryRequest) (*pb.GetAnalysisHistoryResponse, error) {
	analyses, err := h.aiUsecase.GetAnalysisHistory(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get analysis history")
	}

	var analysisIDs []string
	for _, analysis := range analyses {
		analysisIDs = append(analysisIDs, analysis.ID)
	}

	return &pb.GetAnalysisHistoryResponse{AnalysisIds: analysisIDs}, nil
}

func (h *AIGRPCHandler) GetRoadmapByUserID(ctx context.Context, req *pb.GetRoadmapByUserIDRequest) (*pb.GetRoadmapByUserIDResponse, error) {
	roadmap, steps, err := h.aiUsecase.GetRoadmapByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "roadmap not found")
	}

	var pbSteps []*pb.RoadmapStep
	for _, step := range steps {
		pbSteps = append(pbSteps, &pb.RoadmapStep{
			Title:       step.Title,
			Description: step.Description,
			MonthNumber: step.MonthNumber,
			Category:    step.Category,
		})
	}

	return &pb.GetRoadmapByUserIDResponse{
		RoadmapId: roadmap.ID,
		Steps:     pbSteps,
	}, nil
}

func (h *AIGRPCHandler) GetEssayReviewsByUserID(ctx context.Context, req *pb.GetEssayReviewsByUserIDRequest) (*pb.GetEssayReviewsByUserIDResponse, error) {
	reviews, err := h.aiUsecase.GetEssayReviewsByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get reviews")
	}

	var reviewIDs []string
	for _, r := range reviews {
		reviewIDs = append(reviewIDs, r.ID)
	}

	return &pb.GetEssayReviewsByUserIDResponse{ReviewIds: reviewIDs}, nil
}

func (h *AIGRPCHandler) GetEssayReviewByID(ctx context.Context, req *pb.GetEssayReviewByIDRequest) (*pb.GetEssayReviewByIDResponse, error) {
	review, err := h.aiUsecase.GetEssayReviewByID(ctx, req.ReviewId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "review not found")
	}

	return &pb.GetEssayReviewByIDResponse{
		ReviewId:       review.ID,
		EssayText:      review.EssayText,
		GrammarScore:   review.GrammarScore,
		CoherenceScore: review.CoherenceScore,
	}, nil
}

func (h *AIGRPCHandler) DeleteAnalysis(ctx context.Context, req *pb.DeleteAnalysisRequest) (*pb.DeleteAnalysisResponse, error) {
	if err := h.aiUsecase.DeleteAnalysis(ctx, req.AnalysisId); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete analysis")
	}
	return &pb.DeleteAnalysisResponse{Success: true}, nil
}

func (h *AIGRPCHandler) CompareUniversities(ctx context.Context, req *pb.CompareUniversitiesRequest) (*pb.CompareUniversitiesResponse, error) {
	result, chanceOne, chanceTwo, err := h.aiUsecase.CompareUniversities(ctx, req.UniversityIdOne, req.UniversityIdTwo)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to compare")
	}
	return &pb.CompareUniversitiesResponse{
		ComparisonResult: result,
		ChanceOne:        chanceOne,
		ChanceTwo:        chanceTwo,
	}, nil
}

func (h *AIGRPCHandler) GetRecommendedUniversities(ctx context.Context, req *pb.GetRecommendedUniversitiesRequest) (*pb.GetRecommendedUniversitiesResponse, error) {
	universityIDs, err := h.aiUsecase.GetRecommendedUniversities(ctx, req.UserId, req.Gpa)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get recommendations")
	}
	return &pb.GetRecommendedUniversitiesResponse{UniversityIds: universityIDs}, nil
}

func (h *AIGRPCHandler) GetAdmissionAnalysisByID(ctx context.Context, req *pb.GetAdmissionAnalysisByIDRequest) (*pb.GetAdmissionAnalysisByIDResponse, error) {
	analysis, err := h.aiUsecase.GetAdmissionAnalysisByID(ctx, req.AnalysisId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "analysis not found")
	}
	return &pb.GetAdmissionAnalysisByIDResponse{
		AnalysisId:                analysis.ID,
		AdmissionChancePercentage: analysis.AdmissionChancePercentage,
		GapAnalysis:               analysis.GapAnalysis,
	}, nil
}
