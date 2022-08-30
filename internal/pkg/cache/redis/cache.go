package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/logger"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
	"go.opencensus.io/trace"
	"strconv"
	"sync"
	"time"
)

type Redis struct {
	mu      sync.RWMutex
	logger  logger.Logger
	client  *redis.Client
	counter map[string]hitMiss
}

type hitMiss struct {
	miss uint
	hit  uint
}

func NewCache(port string, logger logger.Logger) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Redis{
		logger:  logger,
		client:  client,
		counter: make(map[string]hitMiss, 2),
	}
}

func (r *Redis) addMiss(key string) {
	counter, ok := r.counter[key]
	if !ok {
		counter = hitMiss{}
	}
	counter.miss++
	r.counter[key] = counter
}

func (r *Redis) addHit(key string) {
	counter, ok := r.counter[key]
	if !ok {
		counter = hitMiss{}
	}
	counter.hit++
	r.counter[key] = counter
}

func (r *Redis) CounterInfo() {
	for key, hitMis := range r.counter {
		r.logger.Infof("Key: %s, miss: %d, hit: %d", key, hitMis.miss, hitMis.hit)
	}
}

func (r *Redis) GetUserTickets(ctx context.Context, userId uint) ([]models.Ticket, error) {
	ctx, span := trace.StartSpan(ctx, "cache/redis/GetUserTickets")
	defer span.End()

	key := getKeyUserTickets(userId)
	r.logger.Infof("get cache GetUserTickets %s", key)

	var tickets []models.Ticket

	res := r.client.Get(getKeyUserTickets(userId))
	if res.Err() != nil {
		r.addMiss(key)
		if res.Err() == redis.Nil {
			r.logger.Infof("cache GetUserTickets NIL %v", res.Err())
			return nil, res.Err()
		}
		r.logger.Errorf("cache GetUserTickets %v", res.Err())
		return nil, res.Err()
	}

	cacheResult, err := res.Result()
	if err != nil {
		r.logger.Errorf("cache GetUserTickets cacheResult %v", res)
		return nil, res.Err()
	}

	err = json.Unmarshal([]byte(cacheResult), &tickets)
	if err != nil {
		r.logger.Errorf("cache GetUserTickets Unmarshal %v", res.Err())
		return nil, res.Err()
	}

	r.addHit(key)

	return tickets, nil
}

func (r *Redis) SetUserTickets(ctx context.Context, userId uint, tickets []models.Ticket) bool {
	ctx, span := trace.StartSpan(ctx, "cache/redis/SetUserTickets")
	defer span.End()

	key := getKeyUserTickets(userId)

	r.logger.Infof("set cache SetUserTickets key: %s", key)

	json, _ := json.Marshal(tickets)
	res := r.client.Set(key, string(json), 1*time.Hour)

	if res.Err() != nil {
		r.logger.Errorf("cache SetUserTickets %v", res.Err())
		return false
	}

	return true
}

func (r *Redis) ResetUserTickets(ctx context.Context, userId uint) bool {
	ctx, span := trace.StartSpan(ctx, "cache/redis/ResetUserTickets")
	defer span.End()

	r.logger.Infof("reset cache ResetUserTickets %s", getKeyUserTickets(userId))
	res := r.client.Del(getKeyUserTickets(userId))
	if res.Err() != nil {
		r.logger.Errorf("cache ResetUserTickets %v", res.Err())
		return false
	}

	return true
}

func getKeyUserTickets(userId uint) string {
	return "ticket-userId-" + strconv.FormatUint(uint64(userId), 10)
}
