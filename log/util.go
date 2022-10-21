package log

import (
	"fmt"
	"io"
	"os"
)

// ConfigurePersistentLogging adds a log-to-file writer.
// File content is identical to stdout.
func ConfigurePersistentLogging(fn string, isMux bool) error {
	Info("Logs will be made persistent", "file_name", fn)

	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	if isMux {
		g.logger.SetOutput(io.MultiWriter(g.logger.Out, f))
	} else {
		g.logger.SetOutput(f)
	}

	Info("File logging initialized")
	return nil
}

func ShortHex(bytes []byte) string {
	if len(bytes) > 4 {
		bytes = bytes[:4]
	}
	return fmt.Sprintf("%#x", bytes)
}
