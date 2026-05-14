package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/stepup-ai/api-gateway/proto/university"
)

type UniversityHandler struct {
	universityServiceClient pb.UniversityServiceClient
}

func NewUniversityHandler(universityServiceConnection *grpc.ClientConn) *UniversityHandler {
	return &UniversityHandler{
		universityServiceClient: pb.NewUniversityServiceClient(universityServiceConnection),
	}
}

func (h *UniversityHandler) SearchUniversities(ctx *gin.Context) {
	country := ctx.Query("country")
	universityType := ctx.Query("type")

	response, err := h.universityServiceClient.SearchUniversities(ctx, &pb.SearchUniversitiesRequest{
		Country: country,
		Type:    universityType,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search universities"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"universities": response.Universities})
}

func (h *UniversityHandler) GetUniversityDetails(ctx *gin.Context) {
	universityID := ctx.Param("id")

	response, err := h.universityServiceClient.GetUniversityDetails(ctx, &pb.GetUniversityDetailsRequest{
		UniversityId: universityID,
	})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "university not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"university": response.University})
}

func (h *UniversityHandler) SaveUniversity(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var requestBody struct {
		UniversityID string `json:"university_id"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	_, err := h.universityServiceClient.SaveUniversity(ctx, &pb.SaveUniversityRequest{
		UserId:       userID,
		UniversityId: requestBody.UniversityID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save university"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "university saved successfully"})
}

func (h *UniversityHandler) GetSavedUniversities(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	response, err := h.universityServiceClient.GetSavedUniversities(ctx, &pb.GetSavedUniversitiesRequest{
		UserId: userID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get saved universities"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"universities": response.Universities})
}

func (h *UniversityHandler) SearchGrants(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	country := ctx.Query("country")

	response, err := h.universityServiceClient.SearchGrants(ctx, &pb.SearchGrantsRequest{
		UserId:  userID,
		Country: country,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search grants"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"grants": response.Grants})
}

func (h *UniversityHandler) SaveGrant(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var requestBody struct {
		GrantID string `json:"grant_id"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	_, err := h.universityServiceClient.SaveGrant(ctx, &pb.SaveGrantRequest{
		UserId:  userID,
		GrantId: requestBody.GrantID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save grant"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "grant saved successfully"})
}

func (h *UniversityHandler) GetSavedGrants(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	response, err := h.universityServiceClient.GetSavedGrants(ctx, &pb.GetSavedGrantsRequest{
		UserId: userID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get saved grants"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"grants": response.Grants})
}
