package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	apiPkg "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/api"
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
)

const (
	tokenHeader = "Token"
	authPathRPC = "UserAuth"
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

	go runREST(serverAddress, restAddress)
	runGRPCServer(serverAddress)
}

func runGRPCServer(serverAddress string) {
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error server connect tcp: %v", err)
	}
	defer listener.Close()

	newServer := apiPkg.New()

	option := grpc.UnaryInterceptor(AuthInterceptor)
	grpcServer := grpc.NewServer(option)
	pb.RegisterCinemaServer(grpcServer, newServer)

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func runREST(serverAddress, restAddress string) {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcherREST),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterCinemaHandlerFromEndpoint(ctx, mux, serverAddress, opts); err != nil {
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
	paths := strings.Split(info.FullMethod, "/")

	fmt.Printf("info %#v\n", paths)
	fmt.Printf("info %#v\n", info.FullMethod)

	for _, path := range paths {
		if path == authPathRPC {
			return handler(ctx, req)
		}
	}

	fmt.Printf("info %#v\n", info)
	fmt.Printf("info %#v\n", info.FullMethod)

	metaData, ok := metadata.FromIncomingContext(ctx)
	for k, v := range metaData {
		fmt.Printf("\t%v: %v\n", k, v)
	}

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

	return nil, status.Error(codes.PermissionDenied, "Header [token] is invalid. See 'auth' method")
}
