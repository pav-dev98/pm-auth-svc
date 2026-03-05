package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/pav-dev98/pm-auth-svc/config"
	"github.com/pav-dev98/pm-auth-svc/internal/repository"
	"github.com/pav-dev98/pm-auth-svc/internal/server"
	pb "github.com/pav-dev98/pm-auth-svc/proto/auth"
)

func main() {
	// 1. Cargar config
	cfg := config.Load()

	// 2. Conectar DB
	repo, err := repository.NewAuthRepository(cfg.DSN)
	if err != nil {
		log.Fatalf("error conectando a PostgreSQL: %v", err)
	}
	log.Println("✅ PostgreSQL conectado")

	// 3. Crear servidor gRPC
	grpcServer := grpc.NewServer()
	authServer := server.NewAuthServer(repo, cfg)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	// 4. Levantar el servidor
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("error escuchando en puerto %s: %v", cfg.GRPCPort, err)
	}

	log.Printf("🚀 Auth service corriendo en puerto %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error iniciando gRPC server: %v", err)
	}
}