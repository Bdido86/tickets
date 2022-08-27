package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/cnt"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/logger"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository"
	"go.opencensus.io/trace"
	"time"
)

const timeSleep = 5

type consumer struct {
	Deps
}

type Deps struct {
	Logger           logger.Logger
	CinemaRepository repository.Cinema
}

func NewConsumer(d Deps) *consumer {
	return &consumer{
		Deps: d,
	}
}

func (c *consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			c.Logger.Info("consume Done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				c.Logger.Info("consume channel closed")
				return nil
			}

			cnt.IncTotal()

			ctx := context.Background()
			spanContext, err := getSpanContextFromHeaders(msg)
			if err != nil {
				c.Logger.Errorf("empty X-Span-Context: %v", err)
			}

			ctx, _ = trace.StartSpanWithRemoteParent(ctx, "broker/consumer/topicTicketDelete", spanContext)

			switch msg.Topic {
			case topicTicketCreate:
				var createTicket createTicketStruct
				err := json.Unmarshal(msg.Value, &createTicket)
				if err != nil {
					cnt.IncError()
					c.Logger.Errorf("on unmarshall: %v", err)
					continue
				}

				tokenValue, err := getTokenFromHeaders(msg)
				if err != nil {
					cnt.IncError()
					c.Logger.Errorf("on unmarshall: %v", err)
					continue
				}

				c.createTicket(ctx, &createTicket, tokenValue)
			case topicTicketDelete:
				var deleteTicket deleteTicketStruct
				err := json.Unmarshal(msg.Value, &deleteTicket)
				if err != nil {
					cnt.IncError()
					c.Logger.Errorf("on unmarshall: %v", err)
					continue
				}

				tokenValue, err := getTokenFromHeaders(msg)
				if err != nil {
					cnt.IncError()
					c.Logger.Errorf("on unmarshall: %v", err)
					continue
				}

				c.deleteTicket(ctx, &deleteTicket, tokenValue)
			}

			session.MarkMessage(msg, "")
		}
	}
}

func getTokenFromHeaders(msg *sarama.ConsumerMessage) (string, error) {
	for _, header := range msg.Headers {
		if string(header.Key) == "Token" {
			return string(header.Value), nil
		}
	}
	return "", errors.New("empty token in headers")
}

func getSpanContextFromHeaders(msg *sarama.ConsumerMessage) (trace.SpanContext, error) {
	var spanContext trace.SpanContext
	for _, header := range msg.Headers {
		if string(header.Key) == "X-Span-Context" && len(header.Key) > 0 {
			err := json.Unmarshal(header.Key, &spanContext)
			return spanContext, err
		}
	}
	return spanContext, errors.New("empty X-Span-Context in headers")
}

func (c *consumer) Run(ctx context.Context) error {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, "tickets", cfg)
	if err != nil {
		c.Logger.Fatalf("error run consumer: %v", err)
	}

	for {
		if err := client.Consume(ctx, []string{topicTicketCreate, topicTicketDelete}, c); err != nil {
			c.Logger.Errorf("on consume: %v", err)
			time.Sleep(time.Second * timeSleep)
		}
	}
}

func (c *consumer) createTicket(ctx context.Context, createTicket *createTicketStruct, token string) {
	ctx, span := trace.StartSpan(ctx, "kafka/consumer/GetFilmRoom")
	defer span.End()

	userId, err := c.CinemaRepository.GetUserIdByToken(ctx, token)
	if err != nil {
		cnt.IncError()
		c.Logger.Infof("not found user: %v", err)
		return
	}

	if _, err := c.CinemaRepository.CreateTicket(ctx, createTicket.FilmId, createTicket.PlaceId, userId); err != nil {
		cnt.IncError()
		c.Logger.Errorf("error createTicket: %v", err)
		return
	}

	cnt.IncSuccess()
	c.Logger.Info("success createTicket")
}

func (c *consumer) deleteTicket(ctx context.Context, deleteTicket *deleteTicketStruct, token string) {
	ctx, span := trace.StartSpan(ctx, "kafka/consumer/deleteTicket")
	defer span.End()

	userId, err := c.CinemaRepository.GetUserIdByToken(ctx, token)
	if err != nil {
		cnt.IncError()
		c.Logger.Infof("not found user: %v", err)
		return
	}

	if err := c.CinemaRepository.DeleteTicket(ctx, deleteTicket.Id, userId); err != nil {
		cnt.IncError()
		c.Logger.Errorf("error deleteTicket: %v", err)
		return
	}

	cnt.IncSuccess()
	c.Logger.Info("success deleteTicket")
}
