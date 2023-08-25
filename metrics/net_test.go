package metrics

import (
	"context"
	"fmt"
	"testing"

	"github.com/photon-storage/go-common/testing/require"
)

func TestGetNetstatResults(t *testing.T) {
	t.Skip()

	entries, err := getNetstatResults(context.TODO())
	require.NoError(t, err)
	for _, e := range entries {
		fmt.Printf("%v %v/%v\n", e.iface, e.rx, e.tx)
	}
}
