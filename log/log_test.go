package log_test

import (
	"context"
	"testing"
	"time"

	"github.com/photon-storage/go-common/log"
	"github.com/photon-storage/go-common/testing/require"
)

func TestLog(t *testing.T) {
	t.Skip()

	loc, err := time.LoadLocation("US/Eastern")
	require.NoError(t, err)
	_ = loc

	ctx, cancel := context.WithCancel(context.Background())
	log.Init(&log.Options{
		Context:  ctx,
		LogLevel: log.DebugLevel,
		Sync:     false,
		Location: loc,
	})

	log.Debug("Debug message")
	log.Info("Info message")
	log.Warn("Warn message")
	log.Error("Error message")

	cancel()
	log.WaitForDone()
}
