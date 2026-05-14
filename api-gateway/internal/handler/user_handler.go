package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/stepup-ai/api-gateway/proto/user"
)

type UserHandler struct {
	userServiceClient pb.UserServiceClient
}

func NewUserHandler(userServiceConnection *grpc.ClientConn) *UserHandler {
	return &UserHandler{
		userServiceClient: pb.NewUserServiceClient(userServiceConnection),
	}
}

func (h *UserHandler) RegisterUser(ctx *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.userServiceClient.RegisterUser(ctx, &pb.RegisterUserRequest{
		Email:    requestBody.Email,
		Password: requestBody.Password,
		FullName: requestBody.FullName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
	})
}

func (h *UserHandler) LoginUser(ctx *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.userServiceClient.LoginUser(ctx, &pb.LoginUserRequest{
		Email:    requestBody.Email,
		Password: requestBody.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
	})
}

func (h *UserHandler) LogoutUser(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	_, err := h.userServiceClient.LogoutUser(ctx, &pb.LogoutUserRequest{
		UserId: userID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	var requestBody struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.userServiceClient.RefreshToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: requestBody.RefreshToken,
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
	})
}

func (h *UserHandler) GetUserProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	response, err := h.userServiceClient.GetUserProfile(ctx, &pb.GetUserProfileRequest{
		UserId: userID,
	})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":     response.UserId,
		"gpa":         response.Gpa,
		"sat_score":   response.SatScore,
		"ielts_score": response.IeltsScore,
		"country":     response.Country,
	})
}

func (h *UserHandler) UpdateUserProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var requestBody struct {
		GPA        float32 `json:"gpa"`
		SATScore   int32   `json:"sat_score"`
		IELTSScore float32 `json:"ielts_score"`
		Country    string  `json:"country"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	_, err := h.userServiceClient.UpdateUserProfile(ctx, &pb.UpdateUserProfileRequest{
		UserId:     userID,
		Gpa:        requestBody.GPA,
		SatScore:   requestBody.SATScore,
		IeltsScore: requestBody.IELTSScore,
		Country:    requestBody.Country,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}

func (h *UserHandler) SendPasswordResetEmail(ctx *gin.Context) {
	var requestBody struct {
		Email string `json:"email"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	_, err := h.userServiceClient.SendPasswordResetEmail(ctx, &pb.SendPasswordResetEmailRequest{
		Email: requestBody.Email,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send reset email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "reset email sent successfully"})
}

func (h *UserHandler) ResetPassword(ctx *gin.Context) {
	var requestBody struct {
		ResetToken  string `json:"reset_token"`
		NewPassword string `json:"new_password"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	_, err := h.userServiceClient.ResetPassword(ctx, &pb.ResetPasswordRequest{
		ResetToken:  requestBody.ResetToken,
		NewPassword: requestBody.NewPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}
