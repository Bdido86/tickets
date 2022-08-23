package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/broker/kafka"
	postgres "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository/postgres"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := config.GetConfig()
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DbHost(), c.DbPort(), c.DbUser(), c.DbPassword(), c.DbName())

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatalf("Can't connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Ping database error: %v", err)
	}

	go func() {
		http.ListenAndServe("localhost:8020", nil)
	}()

	deps := kafka.Deps{CinemaRepository: postgres.NewRepository(pool)}
	consumer := kafka.NewConsumer(deps)
	consumer.Run(ctx)
}
