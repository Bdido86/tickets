package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	apiPkg "gitlab.ozon.dev/Bdido86/movie-tickets/internal/api"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
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

	swaggerDir = "./third_party/swagger-ui"
)

func main() {
	c := config.GetConfig()
	if len(c.ServerPort()) == 0 {
		log.Fatal("Config error: server port is empty")
	}
	if len(c.RestPort()) == 0 {
		log.Fatal("Config error: rest port is empty")
	}

	serverAddress := ":" + c.ServerPort()
	restAddress := ":" + c.RestPort()

	go runREST(serverAddress, restAddress, c.RequestTimeOutInMilliSecond())
	runGRPCServer(serverAddress)
}

func runGRPCServer(serverAddress string) {
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error server connect tcp: %v", err)
	}
	defer listener.Close()

	option := grpc.UnaryInterceptor(AuthInterceptor)
	grpcServer := grpc.NewServer(option)
	pb.RegisterCinemaServer(grpcServer, apiPkg.NewServer())

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func runREST(serverAddress, restAddress string, requestTimeoutInMilliSecond time.Duration) {
	ctx := context.Background()

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
	if err := pb.RegisterCinemaHandlerFromEndpoint(ctx, rmux, serverAddress, opts); err != nil {
		panic(err)
	}

	log.Println("Serving gRPC-Gateway on " + restAddress)
	log.Fatalln(http.ListenAndServe(restAddress, mux))
}

func headerMatcherREST(key string) (string, bool) {
	switch key {
	case tokenHeader:
		return key, true
	default:
		return key, false
	}
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if isAuthPath(info.FullMethod) {
		return handler(ctx, req)
	}

	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Header [token] is required1")
	}

	tokens := metaData.Get(tokenHeader)
	if len(tokens) == 0 {
		return nil, status.Error(codes.PermissionDenied, "Header [token] is required2")
	}

	for _, token := range tokens {
		if apiPkg.IsValidToken(token) {
			return handler(ctx, req)
		}
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
