package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	apiPkg "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/api"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
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

	grpcServer := grpc.NewServer()
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
	case "Custom":
		return key, true
	default:
		return key, false
	}
}
