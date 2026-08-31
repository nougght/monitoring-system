package util

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrBatcherFull = errors.New("batcher is full")

type (
	Batcher[T any] struct {
		limit    int
		interval time.Duration
		items    []T
		mu       sync.Mutex
		wg       sync.WaitGroup
		cancel   context.CancelFunc

		signal chan struct{}

		handler batchHandler[T]
		onceRun sync.Once
		running atomic.Bool
	}
	batchHandler[T any] func(context.Context, []T) error
)

func NewBatcher[T any](limit int, interval time.Duration, handler batchHandler[T]) *Batcher[T] {

	batcher := &Batcher[T]{
		limit:    limit,
		interval: interval,
		handler:  handler,
		signal:   make(chan struct{}, 1),
	}
	return batcher
}

func (b *Batcher[T]) Run(ctx context.Context) {
	b.onceRun.Do(func() {
		batcherCtx, cancel := context.WithCancel(ctx)
		b.cancel = cancel
		b.running.Store(true)
		b.runWorker(batcherCtx)
	})
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
				flushContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				b.flush(flushContext) // proccess last items
				return
			case <-ticker.C:
				b.flush(ctx)
			case <-b.signal:
				b.flush(ctx)
			}
		}
	}()
}

func (b *Batcher[T]) flush(ctx context.Context) {
	b.mu.Lock()
	if len(b.items) == 0 {
		b.mu.Unlock()
		return
	}

	batchToProcess := b.items
	b.items = make([]T, 0, b.limit)
	b.mu.Unlock()

	err := b.handler(ctx, batchToProcess)
	if err != nil {
		// return items to buffer
		b.mu.Lock()
		b.items = append(b.items, batchToProcess...)
		b.mu.Unlock()
	}
}

func (b *Batcher[T]) Stop() {
	if b.running.Load() {
		b.cancel()
		b.wg.Wait()
		b.running.Store(false)
	}
}
