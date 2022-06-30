package observability

import (
	"context"
	"sync"
)

var (
	cleanupWaitGroup *sync.WaitGroup = &sync.WaitGroup{}
)

func WaitForCleanup(ctx context.Context) {
	cleanupDone := make(chan struct{})

	go func() {
		cleanupWaitGroup.Wait()

		close(cleanupDone)
	}()

	select {
	case <-ctx.Done():
		return

	case <-cleanupDone:
		return
	}
}
