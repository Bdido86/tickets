package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	apiGrpcServer "gitlab.ozon.dev/Bdido86/movie-tickets/internal/api/grpc/server"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	postgres "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository/postgres"
	pbApiServer "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"strings"
)

const (
	tokenHeader = "Token"
	authPathRPC = "UserAuth"
)

var depsRepo apiGrpcServer.Deps

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := config.GetConfig()

	serverAddress := ":" + c.ServerPort()

	// postgres connection
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DbHost(), c.DbPort(), c.DbUser(), c.DbPassword(), c.DbName())
	// pgBouncer connection
	//psqlConn := fmt.Sprintf("host=%s port=6432 user=%s password=%s dbname=%s sslmode=disable", c.DbHost(), c.DbUser(), c.DbPassword(), c.DbName())

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatalf("Can't connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Ping database error: %v", err)
	}

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error GRPCServer connect tcp: %v", err)
	}
	defer listener.Close()

	depsRepo = apiGrpcServer.Deps{CinemaRepository: postgres.NewRepository(pool)}

	option := grpc.UnaryInterceptor(AuthInterceptor)
	grpcServer := grpc.NewServer(option)
	pbApiServer.RegisterCinemaBackendServer(grpcServer, apiGrpcServer.NewServer(depsRepo))

	log.Println("Serving SERVER GRPC on " + serverAddress)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error GRPCServer listen: %v", err)
	}
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if isAuthPath(info.FullMethod) {
		return handler(ctx, req)
	}

	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Header [token] is required")
	}

	tokens := metaData.Get(tokenHeader)
	if len(tokens) == 0 {
		return nil, status.Error(codes.PermissionDenied, "Header [token] is required")
	}

	for _, token := range tokens {
		userId, err := depsRepo.CinemaRepository.GetUserIdByToken(ctx, token)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "Header [Token] is invalid. See 'auth' method. "+err.Error())
		}

		ctx = context.WithValue(ctx, "userId", userId)
		return handler(ctx, req)
	}

	return nil, status.Error(codes.PermissionDenied, "Header [Token] is invalid. See 'auth' method")
}

func isAuthPath(fullMethod string) bool {
	paths := strings.Split(fullMethod, "/")
	for _, path := range paths {
		if path == authPathRPC {
			return true
		}
	}
	return false
}
