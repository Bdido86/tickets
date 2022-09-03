package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	apiGrpcServer "gitlab.ozon.dev/Bdido86/movie-tickets/internal/api/grpc/server"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	cache "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/cache/redis"
	log "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/logger"
	postgres "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository/postgres"
	pbApiServer "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/server"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"strings"
)

const (
	tokenHeader = "Token"
	authPathRPC = "UserAuth"
)

var depsRepo apiGrpcServer.Deps
var logger log.Logger

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := config.GetConfig()
	logger = log.GetLogger(c.DebugLevel())

	serverAddress := ":" + c.ServerPort()

	// postgres connection
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DbHost(), c.DbPort(), c.DbUser(), c.DbPassword(), c.DbName())
	// pgBouncer connection
	//psqlConn := fmt.Sprintf("host=%s port=6432 user=%s password=%s dbname=%s sslmode=disable", c.DbHost(), c.DbUser(), c.DbPassword(), c.DbName())

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		logger.Fatalf("Can't connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatalf("Ping database error: %v", err)
	}

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		logger.Fatalf("Error GRPCServer connect tcp: %v", err)
	}
	defer listener.Close()

	cache := cache.NewCache(c.RedisAddr(), c.RedisPassword(), c.RedisDb(), logger)
	depsRepo = apiGrpcServer.Deps{
		Logger:           logger,
		CinemaRepository: postgres.NewRepository(pool, logger, cache),
	}

	option := grpc.ChainUnaryInterceptor(
		xSpanInterceptor,
		authInterceptor,
	)
	grpcServer := grpc.NewServer(option)
	pbApiServer.RegisterCinemaBackendServer(grpcServer, apiGrpcServer.NewServer(depsRepo))

	logger.Infof("Serving SERVER GRPC on %s", serverAddress)
	if err = grpcServer.Serve(listener); err != nil {
		logger.Fatalf("Error GRPCServer listen: %v", err)
	}
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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

func xSpanInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	spanContextJson := metaData.Get("X-Span-Context")
	if len(spanContextJson) == 0 {
		logger.Warning("Empty X-Span-Context in header")
		return handler(ctx, req)
	}

	var spanContext trace.SpanContext
	err := json.Unmarshal([]byte(spanContextJson[0]), &spanContext)
	if err != nil {
		logger.Errorf("error  Unmarshal %v", err)
		return handler(ctx, req)
	}

	ctx, span := trace.StartSpanWithRemoteParent(ctx, "gRPC server XSpanInterceptor", spanContext)
	defer span.End()

	return handler(ctx, req)
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
