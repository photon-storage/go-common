package metrics

import (
	"testing"

	"github.com/photon-storage/go-common/testing/require"
)

func TestParseName(t *testing.T) {
	name, labels := parseName("metric_name")
	require.Equal(t, "metric_name", name)
	require.Equal(t, 0, len(labels))

	name, labels = parseName("metric_name.")
	require.Equal(t, "metric_name", name)
	require.Equal(t, 0, len(labels))

	name, labels = parseName("metric_name.label0#value0")
	require.Equal(t, "metric_name", name)
	require.Equal(t, 1, len(labels))
	require.Equal(t, "value0", labels["label0"])

	name, labels = parseName("metric_name.label0#value0.label1#value1")
	require.Equal(t, "metric_name", name)
	require.Equal(t, 2, len(labels))
	require.Equal(t, "value0", labels["label0"])
	require.Equal(t, "value1", labels["label1"])

	name, labels = parseName("metric_name.label0#value0.label1.label2#value2")
	require.Equal(t, "metric_name", name)
	require.Equal(t, 2, len(labels))
	require.Equal(t, "value0", labels["label0"])
	require.Equal(t, "value2", labels["label2"])
}
