package eventbus

import "context"

type EventBus interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Subscribe(ctx context.Context, topic string, handler func([]byte) error) error
	Close() error
}
