package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stepup-ai/university-service/internal/usecase"
	pb "github.com/stepup-ai/university-service/proto/university"
)

type UniversityGRPCHandler struct {
	pb.UnimplementedUniversityServiceServer
	universityUsecase usecase.UniversityUsecase
}

func NewUniversityGRPCHandler(u usecase.UniversityUsecase) *UniversityGRPCHandler {
	return &UniversityGRPCHandler{universityUsecase: u}
}

func (h *UniversityGRPCHandler) SearchUniversities(ctx context.Context, req *pb.SearchUniversitiesRequest) (*pb.SearchUniversitiesResponse, error) {
	universities, err := h.universityUsecase.SearchUniversities(ctx, req.Country, req.Type, req.MinAcceptanceRate, req.MaxAcceptanceRate)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to search universities")
	}

	var pbUniversities []*pb.University
	for _, u := range universities {
		pbUniversities = append(pbUniversities, &pb.University{
			Id:             u.ID,
			Name:           u.Name,
			Country:        u.Country,
			AcceptanceRate: u.AcceptanceRate,
			Type:           u.Type,
			Category:       u.Category,
		})
	}

	return &pb.SearchUniversitiesResponse{Universities: pbUniversities}, nil
}

func (h *UniversityGRPCHandler) GetUniversityDetails(ctx context.Context, req *pb.GetUniversityDetailsRequest) (*pb.GetUniversityDetailsResponse, error) {
	university, err := h.universityUsecase.GetUniversityDetails(ctx, req.UniversityId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "university not found")
	}

	return &pb.GetUniversityDetailsResponse{
		University: &pb.University{
			Id:             university.ID,
			Name:           university.Name,
			Country:        university.Country,
			AcceptanceRate: university.AcceptanceRate,
			Type:           university.Type,
			Category:       university.Category,
		},
	}, nil
}

func (h *UniversityGRPCHandler) SaveUniversity(ctx context.Context, req *pb.SaveUniversityRequest) (*pb.SaveUniversityResponse, error) {
	if err := h.universityUsecase.SaveUniversity(ctx, req.UserId, req.UniversityId); err != nil {
		return nil, status.Error(codes.Internal, "failed to save university")
	}
	return &pb.SaveUniversityResponse{Success: true}, nil
}

func (h *UniversityGRPCHandler) GetSavedUniversities(ctx context.Context, req *pb.GetSavedUniversitiesRequest) (*pb.GetSavedUniversitiesResponse, error) {
	universities, err := h.universityUsecase.GetSavedUniversities(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get saved universities")
	}

	var pbUniversities []*pb.University
	for _, u := range universities {
		pbUniversities = append(pbUniversities, &pb.University{
			Id:             u.ID,
			Name:           u.Name,
			Country:        u.Country,
			AcceptanceRate: u.AcceptanceRate,
			Type:           u.Type,
			Category:       u.Category,
		})
	}

	return &pb.GetSavedUniversitiesResponse{Universities: pbUniversities}, nil
}

func (h *UniversityGRPCHandler) SearchGrants(ctx context.Context, req *pb.SearchGrantsRequest) (*pb.SearchGrantsResponse, error) {
	grants, err := h.universityUsecase.SearchGrants(ctx, req.Country, req.Gpa)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to search grants")
	}

	var pbGrants []*pb.Grant
	for _, g := range grants {
		pbGrants = append(pbGrants, &pb.Grant{
			Id:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Amount:      g.Amount,
			Deadline:    g.Deadline,
		})
	}

	return &pb.SearchGrantsResponse{Grants: pbGrants}, nil
}

func (h *UniversityGRPCHandler) SaveGrant(ctx context.Context, req *pb.SaveGrantRequest) (*pb.SaveGrantResponse, error) {
	if err := h.universityUsecase.SaveGrant(ctx, req.UserId, req.GrantId); err != nil {
		return nil, status.Error(codes.Internal, "failed to save grant")
	}
	return &pb.SaveGrantResponse{Success: true}, nil
}

func (h *UniversityGRPCHandler) GetSavedGrants(ctx context.Context, req *pb.GetSavedGrantsRequest) (*pb.GetSavedGrantsResponse, error) {
	grants, err := h.universityUsecase.GetSavedGrants(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get saved grants")
	}

	var pbGrants []*pb.Grant
	for _, g := range grants {
		pbGrants = append(pbGrants, &pb.Grant{
			Id:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Amount:      g.Amount,
			Deadline:    g.Deadline,
		})
	}

	return &pb.GetSavedGrantsResponse{Grants: pbGrants}, nil
}

func (h *UniversityGRPCHandler) RemoveSavedUniversity(ctx context.Context, req *pb.RemoveSavedUniversityRequest) (*pb.RemoveSavedUniversityResponse, error) {
	if err := h.universityUsecase.RemoveSavedUniversity(ctx, req.UserId, req.UniversityId); err != nil {
		return nil, status.Error(codes.Internal, "failed to remove saved university")
	}
	return &pb.RemoveSavedUniversityResponse{Success: true}, nil
}

func (h *UniversityGRPCHandler) RemoveSavedGrant(ctx context.Context, req *pb.RemoveSavedGrantRequest) (*pb.RemoveSavedGrantResponse, error) {
	if err := h.universityUsecase.RemoveSavedGrant(ctx, req.UserId, req.GrantId); err != nil {
		return nil, status.Error(codes.Internal, "failed to remove saved grant")
	}
	return &pb.RemoveSavedGrantResponse{Success: true}, nil
}

func (h *UniversityGRPCHandler) GetGrantDetails(ctx context.Context, req *pb.GetGrantDetailsRequest) (*pb.GetGrantDetailsResponse, error) {
	grant, err := h.universityUsecase.GetGrantDetails(ctx, req.GrantId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "grant not found")
	}
	return &pb.GetGrantDetailsResponse{
		Grant: &pb.Grant{
			Id:          grant.ID,
			Name:        grant.Name,
			Description: grant.Description,
			Amount:      grant.Amount,
			Deadline:    grant.Deadline,
		},
	}, nil
}

func (h *UniversityGRPCHandler) ListUniversitiesByCategory(ctx context.Context, req *pb.ListUniversitiesByCategoryRequest) (*pb.ListUniversitiesByCategoryResponse, error) {
	universities, err := h.universityUsecase.ListUniversitiesByCategory(ctx, req.Category)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list universities")
	}

	var pbUniversities []*pb.University
	for _, u := range universities {
		pbUniversities = append(pbUniversities, &pb.University{
			Id:             u.ID,
			Name:           u.Name,
			Country:        u.Country,
			AcceptanceRate: u.AcceptanceRate,
			Type:           u.Type,
			Category:       u.Category,
		})
	}

	return &pb.ListUniversitiesByCategoryResponse{Universities: pbUniversities}, nil
}

func (h *UniversityGRPCHandler) GetUniversityStatistics(ctx context.Context, req *pb.GetUniversityStatisticsRequest) (*pb.GetUniversityStatisticsResponse, error) {
	total, reach, target, safety, err := h.universityUsecase.GetUniversityStatistics(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get statistics")
	}
	return &pb.GetUniversityStatisticsResponse{
		TotalUniversities: total,
		TotalReach:        reach,
		TotalTarget:       target,
		TotalSafety:       safety,
	}, nil
}
