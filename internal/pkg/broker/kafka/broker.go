package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/broker"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/logger"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/metadata"
	"log"
)

type Kafka struct {
	asyncProducer sarama.AsyncProducer
	logger        logger.Logger
}

func NewBroker(logger logger.Logger) broker.Broker {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	asyncProducer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		log.Fatalf("asyn kafka: %v", err)
	}

	k := &Kafka{
		asyncProducer: asyncProducer,
		logger:        logger,
	}

	k.checkAsyncProducer()

	return k
}

func (k *Kafka) Close() error {
	defer func() {
		if err := k.asyncProducer.Close(); err != nil {
			k.logger.Fatalln(err)
		}
	}()

	return nil
}

func (k *Kafka) DeleteTicket(ctx context.Context, ticketId uint) error {
	ctx, span := trace.StartSpan(ctx, "broker/DeleteTicket")
	defer span.End()

	request := deleteTicketStruct{
		Id: ticketId,
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		k.logger.Errorf("kafka DeleteTicket %v", err)
		return errors.Wrap(err, "broker/kafka DeleteTicket")
	}

	spanContext, err := getSpanContext(span)
	if err != nil {
		k.logger.Error("empty span context")
	}

	k.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topicTicketDelete,
		Key:   sarama.StringEncoder(topicTicketDelete),
		Value: sarama.ByteEncoder(jsonRequest),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("Token"),
				Value: []byte(getTokenFromContext(ctx)),
			},
			{
				Key:   []byte("X-Span-Context"),
				Value: spanContext,
			},
		},
	}
	return nil
}

func (k *Kafka) CreateTicket(ctx context.Context, filmId, placeId uint) error {
	ctx, span := trace.StartSpan(ctx, "broker/CreateTicket")
	defer span.End()

	request := createTicketStruct{
		FilmId:  filmId,
		PlaceId: placeId,
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		k.logger.Errorf("kafka CreateTicket %v", err)
		return errors.Wrap(err, "broker/kafka CreateTicket")
	}

	spanContext, err := getSpanContext(span)
	if err != nil {
		k.logger.Error("empty span context")
	}

	k.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topicTicketCreate,
		Key:   sarama.StringEncoder(topicTicketCreate),
		Value: sarama.ByteEncoder(jsonRequest),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("Token"),
				Value: []byte(getTokenFromContext(ctx)),
			},
			{
				Key:   []byte("X-Span-Context"),
				Value: spanContext,
			},
		},
	}
	return nil
}

func getTokenFromContext(ctx context.Context) string {
	metaData, _ := metadata.FromIncomingContext(ctx)

	tokens := metaData.Get(token)
	return tokens[0]
}

func getSpanContext(span *trace.Span) ([]byte, error) {
	return json.Marshal(span.SpanContext())
}
