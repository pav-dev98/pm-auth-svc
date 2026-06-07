package grpc

import (
	"context"

	pb "github.com/pav-dev98/pm-proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pav-dev98/pm-auth-svc/internal/application/usecase"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer

	registerUC *usecase.RegisterCredential
	loginUC    *usecase.LoginCredential
	refreshUC  *usecase.RefreshCredential
}

func NewAuthServer(registerUC *usecase.RegisterCredential, loginUC *usecase.LoginCredential,refreshUC *usecase.RefreshCredential) *AuthServer {
	return &AuthServer{
		registerUC: registerUC,
		loginUC:    loginUC,
		refreshUC:  refreshUC,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	accessToken, refreshToken, err := s.registerUC.Execute(req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	accessToken, refreshToken, err := s.loginUC.Execute(req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

func (s *AuthServer) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	accessToken,err := s.refreshUC.Execute(req.RefreshToken)

	if err != nil {
		return nil , status.Error(codes.Unauthenticated,err.Error())
	}
	return &pb.RefreshResponse{
		AccessToken: accessToken,
	},nil
}