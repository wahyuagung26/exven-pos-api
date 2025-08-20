package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventBus interface {
	Publish(ctx context.Context, topic string, event Event) error
	Subscribe(topic string, handler EventHandler) error
	Close() error
}

type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	TenantID  uint64                 `json:"tenant_id"`
	UserID    uint64                 `json:"user_id,omitempty"`
	Data      map[string]interface{} `json:"data"`
}

type EventHandler func(event Event) error

type RabbitMQEventBus struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
	queuePrefix  string
}

func NewRabbitMQEventBus(url, exchangeName, queuePrefix string) (*RabbitMQEventBus, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = channel.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQEventBus{
		conn:         conn,
		channel:      channel,
		exchangeName: exchangeName,
		queuePrefix:  queuePrefix,
	}, nil
}

func (eb *RabbitMQEventBus) Publish(ctx context.Context, topic string, event Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = eb.channel.PublishWithContext(
		ctx,
		eb.exchangeName,
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    event.Timestamp,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (eb *RabbitMQEventBus) Subscribe(topic string, handler EventHandler) error {
	queueName := fmt.Sprintf("%s%s", eb.queuePrefix, topic)

	queue, err := eb.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = eb.channel.QueueBind(
		queue.Name,
		topic,
		eb.exchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	messages, err := eb.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range messages {
			var event Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				msg.Nack(false, false)
				continue
			}

			if err := handler(event); err != nil {
				log.Printf("Failed to handle event: %v", err)
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)
		}
	}()

	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

func (eb *RabbitMQEventBus) Close() error {
	if err := eb.channel.Close(); err != nil {
		return err
	}
	return eb.conn.Close()
}

func NewEvent(eventType string, tenantID, userID uint64, data map[string]interface{}) Event {
	return Event{
		ID:        generateEventID(),
		Type:      eventType,
		Timestamp: time.Now(),
		TenantID:  tenantID,
		UserID:    userID,
		Data:      data,
	}
}

func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}