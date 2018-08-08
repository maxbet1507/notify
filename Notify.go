package notify

import (
	"context"
	"sync"
)

// Notify -
type Notify interface {
	Subscribe() func(context.Context) error
	Broadcast()
}

type rawNotify struct {
	Locker sync.Locker
	m      sync.Mutex
	c      chan struct{}
}

func (s *rawNotify) Subscribe() func(context.Context) error {
	s.m.Lock()
	defer s.m.Unlock()

	c := s.c

	return func(ctx context.Context) error {
		if s.Locker != nil {
			s.Locker.Unlock()
			defer s.Locker.Lock()
		}

		select {
		case <-c:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *rawNotify) Broadcast() {
	s.m.Lock()
	defer s.m.Unlock()

	close(s.c)
	s.c = make(chan struct{})
}

// New -
func New(l sync.Locker) Notify {
	r := &rawNotify{
		Locker: l,
		c:      make(chan struct{}),
	}
	return r
}
