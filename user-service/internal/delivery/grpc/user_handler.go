package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stepup-ai/user-service/internal/usecase"
	pb "github.com/stepup-ai/user-service/proto/user"
)

type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	userUsecase usecase.UserUsecase
}

func NewUserGRPCHandler(u usecase.UserUsecase) *UserGRPCHandler {
	return &UserGRPCHandler{userUsecase: u}
}

func (h *UserGRPCHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	accessToken, refreshToken, err := h.userUsecase.RegisterUser(ctx, req.Email, req.Password, req.FullName)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &pb.RegisterUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *UserGRPCHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	accessToken, refreshToken, err := h.userUsecase.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}
	return &pb.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *UserGRPCHandler) LogoutUser(ctx context.Context, req *pb.LogoutUserRequest) (*pb.LogoutUserResponse, error) {
	if err := h.userUsecase.LogoutUser(ctx, req.UserId); err != nil {
		return nil, status.Error(codes.Internal, "failed to logout")
	}
	return &pb.LogoutUserResponse{Success: true}, nil
}

func (h *UserGRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	newAccessToken, newRefreshToken, err := h.userUsecase.RefreshAccessToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidToken) || errors.Is(err, usecase.ErrTokenExpired) {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}
	return &pb.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (h *UserGRPCHandler) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	profile, err := h.userUsecase.GetUserProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "profile not found")
	}
	return &pb.GetUserProfileResponse{
		UserId:     profile.UserID,
		GPA:        profile.GPA,
		SatScore:   profile.SATScore,
		IeltsScore: profile.IELTSScore,
		Country:    profile.Country,
	}, nil
}

func (h *UserGRPCHandler) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	profile, err := h.userUsecase.GetUserProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "profile not found")
	}
	profile.GPA = req.Gpa
	profile.SATScore = req.SatScore
	profile.IELTSScore = req.IeltsScore
	profile.Country = req.Country
	if err := h.userUsecase.UpdateUserProfile(ctx, profile); err != nil {
		return nil, status.Error(codes.Internal, "failed to update profile")
	}
	return &pb.UpdateUserProfileResponse{Success: true}, nil
}

func (h *UserGRPCHandler) SendPasswordResetEmail(ctx context.Context, req *pb.SendPasswordResetEmailRequest) (*pb.SendPasswordResetEmailResponse, error) {
	if err := h.userUsecase.SendPasswordResetEmail(ctx, req.Email); err != nil {
		return nil, status.Error(codes.Internal, "failed to send reset email")
	}
	return &pb.SendPasswordResetEmailResponse{Success: true}, nil
}

func (h *UserGRPCHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	if err := h.userUsecase.ResetPassword(ctx, req.ResetToken, req.NewPassword); err != nil {
		if errors.Is(err, usecase.ErrResetTokenInvalid) {
			return nil, status.Error(codes.InvalidArgument, "invalid or expired reset token")
		}
		return nil, status.Error(codes.Internal, "failed to reset password")
	}
	return &pb.ResetPasswordResponse{Success: true}, nil
}
