package notify

import (
	"context"
	"sync"
)

// Notify -
type Notify interface {
	sync.Locker
	Subscribe() func(context.Context) error
	Broadcast()
}

type rawNotify struct {
	Locker sync.Locker
	m      sync.Mutex
	c      chan struct{}
}

func (s *rawNotify) Lock() {
	if s.Locker != nil {
		s.Locker.Lock()
	}
}

func (s *rawNotify) Unlock() {
	if s.Locker != nil {
		s.Locker.Unlock()
	}
}

func (s *rawNotify) Subscribe() func(context.Context) error {
	s.m.Lock()
	defer s.m.Unlock()

	c := s.c

	return func(ctx context.Context) error {
		s.Unlock()
		defer s.Lock()

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
