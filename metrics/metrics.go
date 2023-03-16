package metrics

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/photon-storage/go-common/log"
)

var (
	metricNamespace = ""
	counters        = map[string]prometheus.Counter{}
	gauges          = map[string]prometheus.Gauge{}
	histograms      = map[string]prometheus.Histogram{}
)

func Init(ctx context.Context, namespace string, port int) {
	metricNamespace = namespace
	srv := &http.Server{Addr: fmt.Sprintf(":%v", port)}
	http.Handle("/metrics", promhttp.Handler())

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Error("monitoring http ListenAndServe error", "error", err)
		}
	}()

	go func() {
		for range ctx.Done() {
			break
		}
		srv.Close()
	}()
}

// NewCounter declares a new counter.
func NewCounter(name string) {
	metricName, labels := parseName(name)
	counters[name] = promauto.NewCounter(prometheus.CounterOpts{
		Namespace:   metricNamespace,
		Name:        metricName,
		ConstLabels: labels,
	})
}

func CounterInc(name string) {
	c := counters[name]
	if c != nil {
		c.Inc()
	}
}

func CounterAdd(name string, v float64) {
	c := counters[name]
	if c != nil {
		c.Add(v)
	}
}

// NewGauge declares a new gauge.
func NewGauge(name string) {
	metricName, labels := parseName(name)
	gauges[name] = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        metricName,
		ConstLabels: labels,
	})
}

func GaugeInc(name string) {
	g := gauges[name]
	if g != nil {
		g.Inc()
	}
}

func GaugeDec(name string) {
	g := gauges[name]
	if g != nil {
		g.Dec()
	}
}

func GaugeAdd(name string, v float64) {
	g := gauges[name]
	if g != nil {
		g.Add(v)
	}
}

func GaugeSet(name string, v float64) {
	g := gauges[name]
	if g != nil {
		g.Set(v)
	}
}

// NewHistogram declares a new histogram.
// Buckets defines the buckets into which observations are counted. Each
// element in the slice is the upper inclusive bound of a bucket. The
// values must be sorted in strictly increasing order. There is no need
// to add a highest bucket with +Inf bound, it will be added
// implicitly. If Buckets is left as nil or set to a slice of length
// zero, it is replaced by default buckets. The default buckets are
// DefBuckets if no buckets for a native histogram (see below) are used,
// otherwise the default is no buckets. (In other words, if you want to
// use both reguler buckets and buckets for a native histogram, you have
// to define the regular buckets here explicitly.)
func NewHistogram(
	name string,
	buckets ...float64,
) {
	histograms[name] = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: metricNamespace,
		Name:      name,
		Buckets:   buckets,
	})
}

type number interface {
	int | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64
}

func HistAdd[T number](name string, v T) {
	h := histograms[name]
	if h != nil {
		h.Observe(float64(v))
	}
}

func parseName(name string) (string, prometheus.Labels) {
	parts := strings.Split(name, ".")
	if len(parts) == 1 {
		return parts[0], nil
	}

	labels := prometheus.Labels{}
	for i := 1; i < len(parts); i++ {
		kv := strings.Split(parts[i], "#")
		if len(kv) != 2 {
			continue
		}
		labels[kv[0]] = kv[1]
	}

	return parts[0], labels
}
