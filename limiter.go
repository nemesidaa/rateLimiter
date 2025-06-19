package limiter

import (
	"time"

	entrylist "github.com/nemesidaa/rateLimiter/entryList"
)

type Limiter struct {
	data   entrylist.EntryListShape
	ticker *time.Ticker
	stop   chan struct{}
}

func New(limit uint16, period time.Duration) *Limiter {
	return &Limiter{
		data:   entrylist.NewEntryList(limit),
		ticker: time.NewTicker(period),
		stop:   make(chan struct{}),
	}
}

func (lim *Limiter) Increment(entry string) error {
	return lim.data.IncrementFor(entry)
}

func (lim *Limiter) reset() error {
	return lim.data.Reset()
}

func (lim *Limiter) Run() error {
	for {
		select {
		case <-lim.ticker.C:
			if err := lim.reset(); err != nil {
				return err
			}
		case <-lim.stop:
			return nil
		}
	}
}

func (lim *Limiter) Kill() {
	lim.stop <- struct{}{}
	close(lim.stop)
}
