//go:build integration
// +build integration

package tests

import (
	postgresRepo "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository/postgres"
	pbApiServer "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/server"
	"gitlab.ozon.dev/Bdido86/movie-tickets/tests/config"
	"gitlab.ozon.dev/Bdido86/movie-tickets/tests/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var (
	ServerClient pbApiServer.CinemaBackendClient
	Db           *postgres.TDB
	Repository   *postgresRepo.Repository
)

func init() {
	config := config.GetConfig()

	serverAddress := ":" + config.ServerPort()
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTimeout(3*time.Second))
	if err != nil {
		panic(err)
	}

	ServerClient = pbApiServer.NewCinemaBackendClient(conn)
	Db = postgres.NewFromEnv()
	Repository = postgresRepo.NewRepository(Db.Pool)
}
