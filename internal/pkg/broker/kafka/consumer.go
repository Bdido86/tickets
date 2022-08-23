package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/cnt"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository"
	"log"
	"time"
)

const timeSleep = 5

type consumer struct {
	Deps
}

type Deps struct {
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
			log.Print("Done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Print("Data channel closed")
				return nil
			}

			cnt.IncTotal()

			switch msg.Topic {
			case topicTicketCreate:
				var createTicket createTicketStruct
				err := json.Unmarshal(msg.Value, &createTicket)
				if err != nil {
					cnt.IncError()
					log.Printf("on unmarshall: %v", err)
					continue
				}
				c.createTicket(&createTicket)
			case topicTicketDelete:
				var deleteTicket deleteTicketStruct
				err := json.Unmarshal(msg.Value, &deleteTicket)
				if err != nil {
					cnt.IncError()
					log.Printf("on unmarshall: %v", err)
					continue
				}
				c.deleteTicket(&deleteTicket)
			}

			log.Printf("topic: %v, data: %v", msg.Topic, string(msg.Value))
			session.MarkMessage(msg, "")
		}
	}
}

func (c *consumer) Run(ctx context.Context) error {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, "tickets", cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	for {
		if err := client.Consume(ctx, []string{topicTicketCreate, topicTicketDelete}, c); err != nil {
			log.Printf("on consume: %v", err)
			time.Sleep(time.Second * timeSleep)
		}
	}
}

func (c *consumer) createTicket(createTicket *createTicketStruct) {
	ctx := context.Background()
	userId, err := c.CinemaRepository.GetUserIdByToken(ctx, createTicket.Token)
	if err != nil {
		cnt.IncError()
		log.Printf("not Found user: %v", err)
		return
	}

	if _, err := c.CinemaRepository.CreateTicket(ctx, createTicket.FilmId, createTicket.PlaceId, userId); err != nil {
		cnt.IncError()
		log.Printf("error createTicket: %v", err)
		return
	}

	cnt.IncSuccess()
	log.Printf("success CreateTicket")
}

func (c *consumer) deleteTicket(deleteTicket *deleteTicketStruct) {
	ctx := context.Background()
	userId, err := c.CinemaRepository.GetUserIdByToken(ctx, deleteTicket.Token)
	if err != nil {
		cnt.IncError()
		log.Printf("not Found user: %v", err)
		return
	}

	if err := c.CinemaRepository.DeleteTicket(ctx, deleteTicket.Id, userId); err != nil {
		cnt.IncError()
		log.Printf("error deleteTicket: %v", err)
		return
	}

	cnt.IncSuccess()
	log.Printf("success deleteTicket")
}
