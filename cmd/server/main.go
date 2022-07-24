package main

import (
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	apiPkg "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/api"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	c := config.GetConfig()
	if len(c.Port()) == 0 {
		log.Fatal("Config error: port is empty")
	}

	address := ":" + c.Port()
	runGRPCServer(address)
}

func runGRPCServer(address string) {
	listener, err := net.Listen("tcp", address)
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
