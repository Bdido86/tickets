package main

import (
	apiPkg "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/api"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
	"google.golang.org/grpc"
	"log"
	"net"
)

func runGRPCServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error server connect tcp: %v", err)
	}

	newServer := apiPkg.New()

	grpcServer := grpc.NewServer()
	pb.RegisterCinemaServer(grpcServer, newServer)

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

//
//func runREST() {
//	ctx := context.Background()
//
//	ctx, cancel := context.WithCancel(ctx)
//	defer cancel()
//
//	mux := runtime.NewServeMux(
//		runtime.WithIncomingHeaderMatcher(headerMatcherREST),
//	)
//	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
//	if err := pb.RegisterAdminHandlerFromEndpoint(ctx, mux, ":8081", opts); err != nil {
//		panic(err)
//	}
//
//	if err := http.ListenAndServe(":8080", mux); err != nil {
//		panic(err)
//	}
//}
//
//func headerMatcherREST(key string) (string, bool) {
//	switch key {
//	case "Custom":
//		return key, true
//	default:
//		return key, false
//	}
//}
