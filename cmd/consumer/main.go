package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/broker/kafka"
	cache "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/cache/redis"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/logger"
	postgres "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository/postgres"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := config.GetConfig()
	logger := logger.GetLogger(c.DebugLevel())

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DbHost(), c.DbPort(), c.DbUser(), c.DbPassword(), c.DbName())
	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		logger.Fatalf("can't connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatalf("Ping database error: %v", err)
	}

	go func() {
		logger.Info("Start kafka consumer ON localhost:8020")
		http.ListenAndServe("localhost:8020", nil)
	}()

	cache := cache.NewCache(c.RedisAddr(), c.RedisPassword(), c.RedisDb(), logger)
	deps := kafka.Deps{
		Logger:           logger,
		CinemaRepository: postgres.NewRepository(pool, logger, cache),
	}
	consumer := kafka.NewConsumer(deps)
	consumer.Run(ctx)
}
