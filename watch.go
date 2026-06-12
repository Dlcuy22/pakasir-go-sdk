package pakasir

import (
	"context"
	"fmt"
	"time"
)

type WatchHandle struct {
	stop context.CancelFunc
	done chan struct{}
}

func (w *WatchHandle) Stop() {
	w.stop()
	<-w.done
}

func (c *Client) WatchPayment(ctx context.Context, orderID string, amount int, opts WatchOptions) (*WatchHandle, error) {
	orderID = sanitizeUrlSafe(orderID)
	watchKey := fmt.Sprintf("%s_%d", orderID, amount)

	interval := 3 * time.Second
	timeout := 10 * time.Minute

	if opts.Interval > 0 {
		interval = opts.Interval
	}
	if opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	watchCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})

	if timeout > 0 {
		go func() {
			select {
			case <-time.After(timeout):
				cancel()
			case <-watchCtx.Done():
			}
		}()
	}

	ticker := time.NewTicker(interval)
	lastStatus := ""

	go func() {
		defer ticker.Stop()
		defer close(done)
		defer c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.watchers, watchKey)

		check := func() {
			payment, err := c.DetailPayment(watchCtx, orderID, amount)
			if err != nil {
				if opts.OnError != nil {
					opts.OnError(err)
				}
				return
			}

			if lastStatus != payment.Status {
				if opts.OnStatusChange != nil {
					opts.OnStatusChange(payment)
				}
				lastStatus = payment.Status
			}
		}

		check()

		for {
			select {
			case <-ticker.C:
				check()
			case <-watchCtx.Done():
				return
			}
		}
	}()

	c.mu.Lock()
	c.watchers[watchKey] = cancel
	c.mu.Unlock()

	return &WatchHandle{stop: cancel, done: done}, nil
}
