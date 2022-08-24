package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/broker"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/logger"
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

func (k *Kafka) Close(_ context.Context) error {
	defer func() {
		if err := k.asyncProducer.Close(); err != nil {
			k.logger.Fatalln(err)
		}
	}()

	return nil
}

func (k *Kafka) DeleteTicket(ctx context.Context, ticketId uint) error {
	request := deleteTicketStruct{
		Id:    ticketId,
		Token: getToken(ctx),
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		k.logger.Errorf("kafka DeleteTicket %v", err)
		return errors.Wrap(err, "broker/kafka DeleteTicket")
	}

	k.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topicTicketDelete,
		Key:   sarama.StringEncoder(topicTicketDelete),
		Value: sarama.ByteEncoder(jsonRequest),
	}
	return nil
}

func (k *Kafka) CreateTicket(ctx context.Context, filmId, placeId uint) error {
	request := createTicketStruct{
		FilmId:  filmId,
		PlaceId: placeId,
		Token:   getToken(ctx),
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		k.logger.Errorf("kafka CreateTicket %v", err)
		return errors.Wrap(err, "broker/kafka CreateTicket")
	}

	k.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topicTicketCreate,
		Key:   sarama.StringEncoder(topicTicketCreate),
		Value: sarama.ByteEncoder(jsonRequest),
	}
	return nil
}

func getToken(ctx context.Context) string {
	metaData, _ := metadata.FromIncomingContext(ctx)

	tokens := metaData.Get("Token")
	return tokens[0]
}
