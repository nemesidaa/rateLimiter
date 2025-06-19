package limiter_test

import (
	"errors"
	"testing"
	"time"

	limiter "github.com/nemesidaa/rateLimiter"
	entrylist "github.com/nemesidaa/rateLimiter/entryList"
)

func TestLimiter(t *testing.T) {
	t.Run("Basic single-threaded operations", func(t *testing.T) {
		limiter := limiter.New(3, time.Second)
		go func() {
			if err := limiter.Run(); err != nil {
				t.Error(err.Error())
			}
		}()
		for i := 0; i < 3; i++ {
			if err := limiter.Increment("127.0.0.1"); err != nil {
				t.Errorf("Unexpected error on increment %d: %v", i, err)
			}
		}
		if err := limiter.Increment("127.0.0.1"); !errors.Is(err, entrylist.ErrPersonalOverflow) {
			t.Errorf("Expected ErrPersonalOverflow, got %v", err)
		}
		time.Sleep(2 * time.Second)
		if err := limiter.Increment("127.0.0.1"); err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		limiter.Kill()
	})
}
