package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpcApiClient "gitlab.ozon.dev/Bdido86/movie-tickets/internal/api/grpc/client"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/broker/kafka"
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
	clientPort := ":" + c.ClientPort()
	clientGRPCPort := ":" + c.ClientGrpcPort()
	serverPort := ":" + c.ServerPort()

	go runRestGrpc(ctx, clientGRPCPort, serverPort)
	runRest(ctx, clientPort, clientGRPCPort, c.RequestTimeOutInMilliSecond())
}

func runRestGrpc(ctx context.Context, clientGRPCPort string, serverPort string) {
	listener, err := net.Listen("tcp", clientGRPCPort)
	if err != nil {
		log.Fatalf("Error clientGRPCPort connect tcp: %v", err)
	}
	defer listener.Close()

	serverConn, err := grpc.Dial(serverPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error serverPort connect tcp: %v", err)
	}
	defer serverConn.Close()

	broker := kafka.NewBroker()
	defer func() {
		broker.Close(ctx)
	}()

	server := pbApiServer.NewCinemaBackendClient(serverConn)
	clientServer := grpcApiClient.Deps{
		Server: server,
		Broker: broker,
	}

	option := grpc.UnaryInterceptor(AuthInterceptor)
	grpcServer := grpc.NewServer(option)
	pbApiClient.RegisterCinemaFrontendServer(grpcServer, grpcApiClient.NewServer(clientServer))

	log.Println("Serving CLIENT API GRPC on " + clientGRPCPort)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error CLIENT API GRPC listen: %v", err)
	}
}

func runRest(ctx context.Context, clientPort string, clientGRPCPort string, requestTimeoutInMilliSecond time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

	log.Println("Serving CLIENT API on " + clientPort)
	log.Fatalln(http.ListenAndServe(clientPort, mux))
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
