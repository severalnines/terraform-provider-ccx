package events

import (
	"context"
	"sync"
)

type Callback[T any] interface {
	OnCalled(ctx context.Context, payload T)
}

// CallbackMap is a map of callbackID => Callback
type CallbackMap[T any] map[int]Callback[T]

// ListenersMap is a map of itemID => CallbackMap
type ListenersMap[T any] map[string]CallbackMap[T]

type Subscribers[T any] struct {
	ls  ListenersMap[T]
	mut sync.RWMutex
	id  int
}

func (s *Subscribers[T]) Listen(ctx context.Context, itemID string, callback Callback[T]) context.CancelFunc {
	s.mut.Lock()
	defer s.mut.Unlock()

	callbackID := s.id
	s.id += 1

	if s.ls == nil {
		s.ls = make(ListenersMap[T])
	}

	if _, ok := s.ls[itemID]; !ok {
		s.ls[itemID] = make(CallbackMap[T])
	}

	s.ls[itemID][callbackID] = callback

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)

	go s.unlisten(ctx, itemID, callbackID)

	return cancel
}

func (s *Subscribers[T]) unlisten(ctx context.Context, itemID string, callbackID int) {
	<-ctx.Done()

	s.mut.Lock()
	defer s.mut.Unlock()

	if s.ls == nil {
		return
	}

	if _, ok := s.ls[itemID]; !ok {
		return
	}

	if s.ls[itemID] == nil {
		return
	}

	delete(s.ls[itemID], callbackID)
}

func (s *Subscribers[T]) Publish(ctx context.Context, itemID string, payload T) {
	s.mut.RLock()
	defer s.mut.RUnlock()

	if s.ls == nil {
		return
	}

	cb, ok := s.ls[itemID]
	if !ok {
		return
	}

	n := len(cb)
	if n == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(n)

	for i := range cb {
		if cb[i] == nil {
			continue
		}

		go func(c Callback[T]) {
			defer wg.Done()
			c.OnCalled(ctx, payload)
		}(cb[i])
	}

	wg.Wait()
}

func (s *Subscribers[T]) PublishSync(ctx context.Context, itemID string, payload T) {
	s.mut.RLock()
	defer s.mut.RUnlock()

	if s.ls == nil {
		return
	}

	cb, ok := s.ls[itemID]
	if !ok {
		return
	}

	n := len(cb)
	if n == 0 {
		return
	}

	for i := range cb {
		if cb[i] == nil {
			continue
		}

		cb[i].OnCalled(ctx, payload)
	}
}
