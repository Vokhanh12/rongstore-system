// internal/infrastructure/eventbus/rabbitmq_event_bus.go
package eventbus

import (
	"context"
	"fmt"
	"log"
	"server/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQEventBus struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewEventBusFromConfig(cfg *config.Config) (*RabbitMQEventBus, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
	)
	return NewRabbitMQEventBus(url)
}

func NewRabbitMQEventBus(url string) (*RabbitMQEventBus, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQEventBus{
		conn:    conn,
		channel: ch,
	}, nil
}

// Publish gửi message lên exchange hoặc queue.
func (r *RabbitMQEventBus) Publish(ctx context.Context, topic string, message []byte) error {
	// Tạo queue nếu chưa có
	_, err := r.channel.QueueDeclare(
		topic, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Gửi message
	err = r.channel.PublishWithContext(
		ctx,
		"",    // exchange ("" = default)
		topic, // routing key (tên queue)
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("[RabbitMQ] Published message to topic=%s", topic)
	return nil
}

// Subscribe lắng nghe message trên queue và xử lý bằng handler.
func (r *RabbitMQEventBus) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	msgs, err := r.channel.Consume(
		topic, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to consume queue: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgs:
				if err := handler(msg.Body); err != nil {
					log.Printf("[RabbitMQ] handler error: %v", err)
				}
			}
		}
	}()
	log.Printf("[RabbitMQ] Subscribed to topic=%s", topic)
	return nil
}

func (r *RabbitMQEventBus) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}
