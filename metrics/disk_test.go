package metrics

import (
	"context"
	"fmt"
	"testing"

	"github.com/photon-storage/go-common/testing/require"
)

func TestGetDfResults(t *testing.T) {
	t.Skip()

	entries, err := getDfResults(context.TODO())
	require.NoError(t, err)
	for _, e := range entries {
		fmt.Printf("%v %v %v/%v\n", e.fs, e.mount, e.used, e.total)
	}
}
