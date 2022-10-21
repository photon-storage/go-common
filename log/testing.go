package log

import (
	"testing"

	"github.com/sirupsen/logrus/hooks/test"

	"github.com/photon-storage/go-common/testing/require"
)

func TestingHook(t *testing.T) *test.Hook {
	// Enforce to use only in tests
	require.True(t, true)

	h := new(test.Hook)
	g.logger.AddHook(h)
	return h
}
