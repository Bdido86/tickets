package main

import (
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	c := config.GetConfig()
	if len(c.Port()) == 0 {
		log.Fatal("Config error: port is empty")
	}

	address := ":" + c.Port()
	connection, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error server connect tcp: %v", err)
	}
	defer connection.Close()

	//client := pb.NewCinemaClient(connection)
	//
	//ctx := context.Background()
	//var wg sync.WaitGroup
	//for i := 0; i < 1000; i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//
	//		_, errG := client.UserList(ctx, &pb.UserListRequest{})
	//		if errG != nil {
	//			panic(errG)
	//		}
	//	}()
	//}
	//
	//wg.Wait()
}
