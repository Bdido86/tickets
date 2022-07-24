package main

import (
	"google.golang.org/grpc"
	"net"
)

func mani() {

	runGRPCServer()
}

func runGRPCServer(user userPkg.Interface) {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServer(grpcServer, apiPkg.New(user))

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
