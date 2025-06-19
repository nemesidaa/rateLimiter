package entrylist_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	entrylist "github.com/nemesidaa/rateLimiter/entryList"
)

func TestEntryList(t *testing.T) {
	t.Run("Basic single-threaded operations", func(t *testing.T) {
		list := entrylist.NewEntryList(3)

		for i := 0; i < 3; i++ {
			if err := list.IncrementFor("127.0.0.1"); err != nil {
				t.Errorf("Unexpected error on increment %d: %v", i, err)
			}
		}

		if err := list.IncrementFor("127.0.0.1"); !errors.Is(err, entrylist.ErrPersonalOverflow) {
			t.Errorf("Expected ErrPersonalOverflow, got %v", err)
		}

		if err := list.Reset(); err != nil {
			t.Errorf("Reset failed: %v", err)
		}

		if err := list.IncrementFor("127.0.0.1"); err != nil {
			t.Errorf("Unexpected error after reset: %v", err)
		}
	})

	t.Run("Concurrent increments", func(t *testing.T) {
		list := entrylist.NewEntryList(100)
		const workers = 50
		const increments = 100
		var wg sync.WaitGroup

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < increments; j++ {
					_ = list.IncrementFor("192.168.1.1")
				}
			}()
		}

		wg.Wait()

	})

	t.Run("Reset during operations", func(t *testing.T) {
		list := entrylist.NewEntryList(1000)
		var wg sync.WaitGroup
		stop := make(chan struct{})

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-stop:
						return
					default:
						err := list.IncrementFor("10.0.0.1")
						if errors.Is(err, entrylist.ErrDownStated) {
							t.Logf("Поймано состояние перезапуска мапы: %s", err.Error())
						}
					}
				}
			}()
		}

		// Периодически сбрасываем
		for i := 0; i < 5; i++ {
			time.Sleep(10 * time.Millisecond)
			if err := list.Reset(); err != nil {
				t.Errorf("Reset failed: %v", err)
			}
		}

		close(stop)
		wg.Wait()

	})

}
