package util

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrBatcherFull = errors.New("batcher is full")

type Batcher[T any] struct {
	limit    int
	interval time.Duration
	items    []T
	mu       sync.Mutex
	wg       sync.WaitGroup
	cancel   context.CancelFunc

	signal chan struct{}

	handler func([]T)
}

func NewBatcher[T any](ctx context.Context, limit int, interval time.Duration, handler func([]T)) *Batcher[T] {
	batcherCtx, cancel := context.WithCancel(ctx)
	batcher := &Batcher[T]{
		limit:    limit,
		interval: interval,
		handler:  handler,
		signal:   make(chan struct{}, 1),
		cancel:   cancel,
	}
	batcher.runWorker(batcherCtx)
	return batcher
}

func (b *Batcher[T]) Add(item T) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.items) >= b.limit {
		return ErrBatcherFull
	}

	b.items = append(b.items, item)
	if len(b.items) >= b.limit {
		select {
		case b.signal <- struct{}{}:
		default:
		}
	}
	return nil
}

func (b *Batcher[T]) runWorker(ctx context.Context) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		ticker := time.NewTicker(b.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				b.flush() // proccess last items
				return
			case <-ticker.C:
				b.flush()
			case <-b.signal:
				b.flush()
			}
		}
	}()
}

func (b *Batcher[T]) flush() {
	b.mu.Lock()
	if len(b.items) == 0 {
		b.mu.Unlock()
		return
	}

	batchToProcess := b.items
	b.items = make([]T, 0, b.limit)
	b.mu.Unlock()

	b.handler(batchToProcess)
}

func (b *Batcher[T]) Stop() {
	b.cancel()
	b.wg.Wait()
}
