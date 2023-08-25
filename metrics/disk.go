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

type dfEntry struct {
	fs    string
	mount string
	total int64
	used  int64
}

func getDfResults(ctx context.Context) ([]*dfEntry, error) {
	var disks []*dfEntry
	switch runtime.GOOS {
	case "linux":
		cmd := exec.Command(
			"/bin/sh",
			"-c",
			`df -k --output=source,target,size,used |
			     grep -v tmpfs |
			     grep -v '/var/lib/docker/overlay' |
				 grep -v 'Filesystem'`,
		)
		cmd.Env = append(os.Environ())
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}

		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) != 4 {
				continue
			}

			totalBlks, err := strconv.ParseInt(fields[2], 10, 64)
			if err != nil {
				return nil, err
			}

			usedBlks, err := strconv.ParseInt(fields[3], 10, 64)
			if err != nil {
				return nil, err
			}

			disks = append(disks, &dfEntry{
				fs:    strings.TrimSpace(fields[0]),
				mount: strings.TrimSpace(fields[1]),
				total: totalBlks * 1024,
				used:  usedBlks * 1024,
			})
		}
	case "darwin":
	case "windows":
	}

	return disks, nil
}

func dfLabel(e *dfEntry) string {
	return fmt.Sprintf("fs#%v.mount#%v", e.fs, e.mount)
}

func RegisterDiskMetrics(ctx context.Context) error {
	disks, err := getDfResults(ctx)
	if err != nil {
		return err
	}

	registeredDisks := map[string]bool{}
	for _, d := range disks {
		lbl := dfLabel(d)
		registeredDisks[lbl] = true
		NewGauge("host_disk_total_bytes." + lbl)
		NewGauge("host_disk_used_bytes." + lbl)
	}

	go diskMetricsUpdateLoop(ctx, registeredDisks)

	return nil
}

func diskMetricsUpdateLoop(
	ctx context.Context,
	registeredDisks map[string]bool,
) {
	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			disks, _ := getDfResults(ctx)

			reported := map[string]bool{}
			for _, d := range disks {
				lbl := dfLabel(d)
				if !registeredDisks[lbl] {
					continue
				}
				GaugeSet(
					"host_disk_total_bytes."+lbl,
					float64(d.total),
				)
				GaugeSet(
					"host_disk_used_bytes."+lbl,
					float64(d.used),
				)
				reported[lbl] = true
			}

			for lbl := range registeredDisks {
				if reported[lbl] {
					continue
				}
				GaugeSet("host_disk_total_bytes."+lbl, 0)
				GaugeSet("host_disk_used_bytes."+lbl, 0)
			}
		}
	}
}
