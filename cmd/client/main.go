package main

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	const defaultName = "defaultName"

	c := config.GetConfig()
	if len(c.ServerPort()) == 0 {
		log.Fatal("Config error: port is empty")
	}

	address := ":" + c.ServerPort()
	connection, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error server connect tcp: %v", err)
	}
	defer connection.Close()

	client := pb.NewCinemaClient(connection)

	ctx := context.Background()

	request := &pb.UserAuthRequest{
		Name: defaultName,
	}

	response, err := client.UserAuth(ctx, request)

	fmt.Printf("%+v \n", response)

}
