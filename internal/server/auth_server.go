package server

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pav-dev98/pm-auth-svc/config"
	"github.com/pav-dev98/pm-auth-svc/internal/repository"
	pb "github.com/pav-dev98/pm-proto/auth"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	repo   *repository.AuthRepository
	config *config.Config
}

func NewAuthServer(repo *repository.AuthRepository, cfg *config.Config) *AuthServer {
	return &AuthServer{repo: repo, config: cfg}
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 1. Buscar usuario en DB
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "error interno")
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "credenciales inválidas")
	}

	// 2. Verificar password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "credenciales inválidas")
	}

	// 3. Generar JWT
	token, err := s.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "error generando token")
	}

	return &pb.LoginResponse{
		AccessToken:  token,
		RefreshToken: "", // Redis lo implementamos después
		TokenType:    "Bearer",
	}, nil
}

func (s *AuthServer) generateJWT(userID uint, email string) (string, error) {
	expiration, _ := time.ParseDuration(s.config.JWTExpiration)

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(expiration).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}