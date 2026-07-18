package event

import "context"

type Event interface {
	Name() string
	Subject() string
}

type EventHandler func(ctx context.Context, event Event)

type EventBus interface {
	Start() error
	Shutdown(ctx context.Context) error
	Subscribe(subject string, handler EventHandler) error
	Publish(event Event)
}
