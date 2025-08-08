package pubsub

import (
    "context"
    "sync"
)

type Event[T any] interface {
    Publish(ctx context.Context, eventArg T)
    Subscribe(func(ctx context.Context, eventArg T))
    Clear() // Clear clears all subscribers.
}

type event[T any] struct {
    sync.RWMutex
    subscribers []func(ctx context.Context, eventArg T)
}

// New returns a thread-safe event that can add subscribers while publishing.
func New[T any]() Event[T] { return &event[T]{} }

func (ts *event[T]) Publish(ctx context.Context, eventArg T) {
    ts.RWMutex.RLock()
    defer ts.RWMutex.RUnlock()

    for _, fn := range ts.subscribers {
        fn(ctx, eventArg)
    }
}

func (ts *event[T]) Subscribe(fn func(ctx context.Context, eventArg T)) {
    ts.RWMutex.Lock()
    defer ts.RWMutex.Unlock()

    ts.subscribers = append(ts.subscribers, fn)
}

func (ts *event[T]) Clear() {
    ts.RWMutex.Lock()
    defer ts.RWMutex.Unlock()

    ts.subscribers = nil
}

type simple[T any] struct {
    subscribers []func(ctx context.Context, eventArg T)
}

// NewSimple returns a simple event that doesn't guarantee concurrency.
func NewSimple[T any]() Event[T] { return &simple[T]{} }

func (s *simple[T]) Publish(ctx context.Context, eventArg T) {
    for _, fn := range s.subscribers {
        fn(ctx, eventArg)
    }
}

func (s *simple[T]) Subscribe(fn func(ctx context.Context, eventArg T)) {
    s.subscribers = append(s.subscribers, fn)
}

func (s *simple[T]) Clear() {
    s.subscribers = nil
}
