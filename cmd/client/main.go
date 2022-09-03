package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpcApiClient "gitlab.ozon.dev/Bdido86/movie-tickets/internal/api/grpc/client"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/broker/kafka"
	cache "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/cache/redis"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/logger"
	pbApiClient "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/client"
	pbApiServer "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	tokenHeader = "Token"
	authPathRPC = "UserAuth"
	swaggerDir  = "./third_party/swagger-ui"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := config.GetConfig()
	logger := logger.GetLogger(c.DebugLevel())

	clientPort := ":" + c.ClientPort()
	clientGRPCPort := ":" + c.ClientGrpcPort()
	serverPort := ":" + c.ServerPort()

	go runRestGrpc(ctx, logger, clientGRPCPort, serverPort, c)
	runRest(ctx, logger, clientPort, clientGRPCPort, c.RequestTimeOutInMilliSecond())
}

func runRestGrpc(_ context.Context, logger logger.Logger, clientGRPCPort string, serverPort string, config *config.Config) {
	listener, err := net.Listen("tcp", clientGRPCPort)
	if err != nil {
		logger.Fatalf("Error clientGRPCPort connect tcp: %v", err)
	}
	defer listener.Close()

	serverConn, err := grpc.Dial(serverPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Error serverPort connect tcp: %v", err)
	}
	defer serverConn.Close()

	broker := kafka.NewBroker(logger)
	defer func() {
		broker.Close()
	}()

	cache := cache.NewCache(config.RedisAddr(), config.RedisPassword(), config.RedisDb(), logger)

	server := pbApiServer.NewCinemaBackendClient(serverConn)
	clientServer := grpcApiClient.Deps{
		Server: server,
		Broker: broker,
		Logger: logger,
		Cache:  cache,
	}

	option := grpc.UnaryInterceptor(AuthInterceptor)
	grpcServer := grpc.NewServer(option)
	pbApiClient.RegisterCinemaFrontendServer(grpcServer, grpcApiClient.NewServer(clientServer))

	logger.Infof("Serving CLIENT API GRPC on %s", clientGRPCPort)
	if err = grpcServer.Serve(listener); err != nil {
		logger.Fatalf("Error CLIENT API GRPC listen: %v", err)
	}
}

func runRest(ctx context.Context, logger logger.Logger, clientPort string, clientGRPCPort string, requestTimeoutInMilliSecond time.Duration) {
	runtime.DefaultContextTimeout = requestTimeoutInMilliSecond
	rmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcherREST),
	)

	// Swagger server
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	fs := http.FileServer(http.Dir(swaggerDir))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pbApiClient.RegisterCinemaFrontendHandlerFromEndpoint(ctx, rmux, clientGRPCPort, opts); err != nil {
		log.Fatalf("Error clientGRPC listen: %v", err)
	}

	logger.Fatalln(http.ListenAndServe(clientPort, mux))
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
	if len(tokens[0]) < 30 {
		return nil, status.Error(codes.PermissionDenied, "Header [token] is invalid")
	}

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

func headerMatcherREST(key string) (string, bool) {
	switch key {
	case tokenHeader:
		return key, true
	default:
		return key, false
	}
}
