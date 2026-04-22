package main
// github.com/pav-dev98/pm-auth-svc
import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/pav-dev98/pm-auth-svc/config"
	pb "github.com/pav-dev98/pm-proto/auth"

	// infrastructure
	"github.com/pav-dev98/pm-auth-svc/internal/infrastructure/persistence/postgress"
	"github.com/pav-dev98/pm-auth-svc/internal/infrastructure/security/bcrypt"
	"github.com/pav-dev98/pm-auth-svc/internal/infrastructure/security/jwt"

	// usecases
	"github.com/pav-dev98/pm-auth-svc/internal/application/usecase"

	// interfaces
	grpcHandler "github.com/pav-dev98/pm-auth-svc/internal/interfaces/grpc"
)

func main() {

	// 1. Config
	cfg := config.Load()

	// 2. Infraestructura
	repo, err := postgress.NewAuthRepository(cfg.DSN)
	if err != nil {
		log.Fatalf("error conectando a PostgreSQL: %v", err)
	}

	hasher := bcrypt.NewBcryptHasher()

	expiration, _ := time.ParseDuration(cfg.JWTExpiration)
	tokenService := jwt.NewJWTService(cfg.JWTSecret, expiration)

	log.Println("✅ Infraestructura inicializada")

	// 3. Use cases
	registerUC := usecase.NewRegisterCredential(repo, hasher, tokenService)
	loginUC := usecase.NewLoginCredential(repo, hasher, tokenService)

	// 4. gRPC handler (adapter)
	authServer := grpcHandler.NewAuthServer(registerUC, loginUC)

	// 5. gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	reflection.Register(grpcServer)

	// 6. Listen
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("error escuchando en puerto %s: %v", cfg.GRPCPort, err)
	}

	log.Printf("🚀 Auth service corriendo en puerto %s", cfg.GRPCPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error iniciando gRPC server: %v", err)
	}
}
