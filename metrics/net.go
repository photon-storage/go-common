package metrics

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type nsEntry struct {
	iface string
	rx    int64
	tx    int64
}

func getNetstatResults(ctx context.Context) ([]*nsEntry, error) {
	var ifaces []*nsEntry
	switch runtime.GOOS {
	case "linux":
		cmd := exec.Command(
			"/bin/sh",
			"-c",
			`netstat -i |
				 grep -v 'Kernel Interface table' |
				 grep -v Iface |
				 grep -v MTU`,
		)
		cmd.Env = append(os.Environ())
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}

		for _, line := range strings.Split(string(out), "\n") {
			fields := strings.Fields(strings.TrimSpace(line))
			if len(fields) != 11 {
				continue
			}

			rx, err := strconv.ParseInt(fields[2], 10, 64)
			if err != nil {
				return nil, err
			}

			tx, err := strconv.ParseInt(fields[6], 10, 64)
			if err != nil {
				return nil, err
			}

			ifaces = append(ifaces, &nsEntry{
				iface: strings.TrimSpace(fields[0]),
				rx:    rx,
				tx:    tx,
			})
		}
	case "darwin":
	case "windows":
	}

	return ifaces, nil
}

func ifaceLabel(e *nsEntry) string {
	return fmt.Sprintf("iface#%v", e.iface)
}

func RegisterIfaceMetrics(ctx context.Context) error {
	ifaces, err := getNetstatResults(ctx)
	if err != nil {
		return err
	}

	registeredIfaces := map[string]bool{}
	for _, iface := range ifaces {
		lbl := ifaceLabel(iface)
		registeredIfaces[lbl] = true
		NewGauge("host_iface_rx_total_bytes." + lbl)
		NewGauge("host_iface_tx_total_bytes." + lbl)
	}

	go ifaceMetricsUpdateLoop(ctx, registeredIfaces)

	return nil
}

func ifaceMetricsUpdateLoop(
	ctx context.Context,
	registeredIfaces map[string]bool,
) {
	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			ifaces, _ := getNetstatResults(ctx)

			reported := map[string]bool{}
			for _, iface := range ifaces {
				lbl := ifaceLabel(iface)
				if !registeredIfaces[lbl] {
					continue
				}
				GaugeSet(
					"host_iface_rx_total_bytes."+lbl,
					float64(iface.rx),
				)
				GaugeSet(
					"host_iface_tx_total_bytes."+lbl,
					float64(iface.tx),
				)
				reported[lbl] = true
			}

			for lbl := range registeredIfaces {
				if reported[lbl] {
					continue
				}
				GaugeSet("host_iface_rx_total_bytes."+lbl, 0)
				GaugeSet("host_iface_tx_total_bytes."+lbl, 0)
			}
		}
	}
}
