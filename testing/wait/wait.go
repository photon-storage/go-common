package wait

import (
	"context"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

// ForCond waits for eval function to return true before ctx is done.
func ForCond(t *testing.T, ctx context.Context, eval func() bool) {
	tk := time.NewTicker(10 * time.Millisecond)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			if eval() {
				return
			}

		case <-ctx.Done():
			_, file, line, _ := runtime.Caller(1)
			t.Errorf(
				"ForCond(): condition NOT met before context done: %s:%d",
				filepath.Base(file),
				line,
			)
		}
	}
}

// ForGroup waits for a WaitGroup to resolve before ctx is done.
func ForGroup(t *testing.T, ctx context.Context, wg *sync.WaitGroup) {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	select {
	case <-ch:
		return

	case <-ctx.Done():
		_, file, line, _ := runtime.Caller(1)
		t.Errorf(
			"ForGroup() condition NOT met before context done: %s:%d",
			filepath.Base(file),
			line,
		)
	}
}

// ForChan waits for channel to receive data before ctx is done.
func ForChan[T any](t *testing.T, ctx context.Context, ch chan T, nilVal T) T {
	select {
	case v := <-ch:
		return v

	case <-ctx.Done():
		_, file, line, _ := runtime.Caller(1)
		t.Errorf(
			"ForChan() condition NOT met before context done: %s:%d",
			filepath.Base(file),
			line,
		)
	}

	return nilVal
}
