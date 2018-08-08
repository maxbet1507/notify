package notify_test

import (
	"context"
	"runtime"
	"sync"
	"testing"

	"github.com/maxbet1507/notify"
)

func TestNotify1(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	n := notify.New(nil)

	wg := sync.WaitGroup{}

	wg.Add(2)
	subs := n.Subscribe()
	go func() {
		defer wg.Done()
		if err := subs(ctx); err != nil {
			t.Fatal(err)
		}
	}()
	go func() {
		defer wg.Done()
		n.Broadcast()
	}()
	wg.Wait()

	wg.Add(2)
	subs = n.Subscribe()
	go func() {
		defer wg.Done()
		if err := subs(ctx); err != ctx.Err() {
			t.Fatal(err)
		}
	}()
	go func() {
		defer wg.Done()
		cancel()
	}()
	wg.Wait()
}

func TestNotify2(t *testing.T) {
	ctx := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(3)

	x := 0
	l := &sync.Mutex{}
	n := notify.New(l)

	go func() {
		defer wg.Done()

		l.Lock()
		defer l.Unlock()

		x = 1
		n.Subscribe()(ctx)
		if x != 2 {
			t.Fatal(x)
		}
		x = 3
		n.Broadcast()
	}()

	go func() {
		defer wg.Done()

		l.Lock()
		defer l.Unlock()

		for {
			switch x {
			case 1:
				x = 2
				n.Broadcast()
				return
			}
			l.Unlock()
			runtime.Gosched()
			l.Lock()
		}
	}()

	go func() {
		defer wg.Done()

		l.Lock()
		defer l.Unlock()

		for {
			switch x {
			case 2:
				n.Subscribe()(ctx)
				if x != 3 {
					t.Fatal(x)
				}
				return
			case 3:
				return
			}
			l.Unlock()
			runtime.Gosched()
			l.Lock()
		}
	}()

	wg.Wait()
}
