package eventbus

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"

	"github.com/nougght/monitoring-system/server/internal/config"
	model "github.com/nougght/monitoring-system/server/internal/model/event"
)

type SubChans []chan model.Event

var (
	ErrInvalidSubjecName       = errors.New("invalid subject name")
	ErrFullSubjectNameRequired = errors.New("full subject name required")
)

type EventBus struct {
	cfg  *config.Config
	mu   sync.RWMutex
	subs map[string]SubChans

	busCtx       context.Context
	candelBusCtx context.CancelFunc
	once         sync.Once
	wg           sync.WaitGroup
}

func NewEventBus(cfg *config.Config) *EventBus {
	return &EventBus{
		cfg:  cfg,
		subs: make(map[string]SubChans),
	}
}

func (b *EventBus) Start(ctx context.Context) {
	b.busCtx, b.candelBusCtx = context.WithCancel(ctx)
}

func (b *EventBus) Shutdown(ctx context.Context) error {
	b.once.Do(func() {
		b.candelBusCtx()
		b.mu.Lock()
		defer b.mu.Unlock()
		for _, subChans := range b.subs {
			for _, ch := range subChans {
				close(ch)
			}
		}
	})

	wgDone := make(chan struct{})

	select {
	case <-wgDone:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *EventBus) Subscribe(subject string, handler model.EventHandler, bufferSize int, droppedCount *int) error {
	if !isValidSubject(subject) {
		return ErrInvalidSubjecName
	}
	if bufferSize == 0 {
		bufferSize = 1
	}
	ch := make(chan model.Event, bufferSize)
	b.subs[subject] = append(b.subs[subject], ch)

	go func() {
		for {
			if b.busCtx == nil {
				continue
			}
			select {
			case <-b.busCtx.Done():
				return
			case event, ok := <-ch:
				if !ok {
					return
				}
				handler(b.busCtx, event)
			}
		}
	}()
	return nil
}

func (b *EventBus) Publish(ctx context.Context, event model.Event) error {
	if !isFullSubject(event.Subject()) {
		return ErrFullSubjectNameRequired
	}
	isWildcard, tokens := isWildcard(event.Subject())
	if !isWildcard {
		b.mu.RLock()
		defer b.mu.RUnlock()
		subs, ok := b.subs[event.Subject()]
		if ok {
			for _, sub := range subs {
				sub <- event
			}
		}
		return nil
	}
	for subject, subs := range b.subs {
		if matchFullSubject(strings.Split(subject, "."), tokens) {
			for _, sub := range subs {
				sub <- event
			}
		}
	}
	return nil
}

func isValidSubject(subject string) bool {
	if i := strings.Index(subject, ">"); i != -1 && i != len(subject)-1 {
		return false
	}
	return true
}

func isFullSubject(subject string) bool {
	if strings.ContainsAny(subject, ">*") {
		return false
	}
	return true
}

func isWildcard(subject string) (bool, []string) {
	tokens := strings.Split(subject, ".")
	if tokens[len(tokens)-1] == ">" {
		return true, tokens
	}
	if slices.Contains(tokens, "*") {
		return true, tokens
	}
	return false, nil
}

func matchFullSubject(subject, fullSubject []string) bool {
	for i, token := range subject {
		if i >= len(subject) || token != "*" && token != fullSubject[i] {
			return false
		}
		if token == ">" {
			return true
		}
	}

	return len(subject) == len(fullSubject)
}
