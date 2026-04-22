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
}

func NewAuthServer(registerUC *usecase.RegisterCredential, loginUC *usecase.LoginCredential) *AuthServer {
	return &AuthServer{
		registerUC: registerUC,
		loginUC:    loginUC,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	token, _, err := s.registerUC.Execute(req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		AccessToken:  token,
		TokenType:    "Bearer",
		RefreshToken: "",
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	token, err := s.loginUC.Execute(req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil
}